// Package schema verifies deployment metadata without mutating the database.
package schema

import (
	"context"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"strings"
)

//go:embed columns.json constraints.csv
var files embed.FS

type Column struct {
	Table      string  `json:"table_name"`
	Name       string  `json:"column_name"`
	Type       string  `json:"udt_name"`
	Nullable   string  `json:"is_nullable"`
	Generated  string  `json:"is_generated"`
	Default    *string `json:"column_default"`
	Expression *string `json:"generation_expression"`
	Length     *int    `json:"character_maximum_length"`
	Precision  *int    `json:"numeric_precision"`
	Scale      *int    `json:"numeric_scale"`
}

func Check(ctx context.Context, pool *pgxpool.Pool) error {
	b, _ := files.ReadFile("columns.json")
	var expected []Column
	if e := json.Unmarshal(b, &expected); e != nil {
		return fmt.Errorf("invalid schema baseline")
	}
	rows, e := pool.Query(ctx, `SELECT table_name,column_name,udt_name,is_nullable,is_generated,column_default,generation_expression,character_maximum_length,numeric_precision,numeric_scale FROM information_schema.columns WHERE table_schema='public'`)
	if e != nil {
		return fmt.Errorf("cannot read schema metadata")
	}
	actual := map[string]Column{}
	for rows.Next() {
		var c Column
		if e = rows.Scan(&c.Table, &c.Name, &c.Type, &c.Nullable, &c.Generated, &c.Default, &c.Expression, &c.Length, &c.Precision, &c.Scale); e != nil {
			rows.Close()
			return fmt.Errorf("cannot decode schema metadata")
		}
		actual[c.Table+"."+c.Name] = c
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return fmt.Errorf("cannot read schema metadata")
	}
	for _, c := range expected {
		if !reflect.DeepEqual(actual[c.Table+"."+c.Name], c) {
			return fmt.Errorf("schema mismatch: %s.%s", c.Table, c.Name)
		}
	}
	constraints, e := pool.Query(ctx, `SELECT t.relname,c.conname,pg_get_constraintdef(c.oid) FROM pg_constraint c JOIN pg_class t ON t.oid=c.conrelid JOIN pg_namespace n ON n.oid=t.relnamespace WHERE n.nspname='public'`)
	if e != nil {
		return fmt.Errorf("cannot read constraints")
	}
	found := map[string]string{}
	for constraints.Next() {
		var table, name, def string
		if e = constraints.Scan(&table, &name, &def); e != nil {
			constraints.Close()
			return fmt.Errorf("cannot decode constraints")
		}
		found[table+"."+name] = def
	}
	e = constraints.Err()
	constraints.Close()
	if e != nil {
		return fmt.Errorf("cannot read constraints")
	}
	b, _ = files.ReadFile("constraints.csv")
	records, e := csv.NewReader(strings.NewReader(string(b))).ReadAll()
	if e != nil {
		return fmt.Errorf("invalid constraint baseline")
	}
	for _, r := range records[1:] {
		if len(r) != 3 || found[r[0]+"."+r[1]] != r[2] {
			return fmt.Errorf("constraint baseline mismatch")
		}
	}
	sets := map[string][]string{"role_enum": {"supervisor", "assistant", "worker"}, "leave_status_enum": {"pending", "approved", "rejected"}, "contract_status_enum": {"registered", "active", "complete", "cancelled"}, "expense_type_enum": {"fund_transfer", "actual_expense"}, "invoice_status_enum": {"pending", "paid"}, "shift_status_enum": {"scheduled", "completed", "cancelled"}, "attendance_status_enum": {"on_time", "late", "absent", "leave"}, "requisition_status_enum": {"pending_survey", "pending_procurement", "pending_fund", "pending_approval", "completed", "rejected"}, "requisition_type_enum": {"tor_base", "additional"}}
	for name, labels := range sets {
		var got []string
		if e := pool.QueryRow(ctx, `SELECT array_agg(e.enumlabel ORDER BY e.enumsortorder) FROM pg_type t JOIN pg_namespace n ON n.oid=t.typnamespace JOIN pg_enum e ON e.enumtypid=t.oid WHERE n.nspname='public' AND t.typname=$1`, name).Scan(&got); e != nil {
			return fmt.Errorf("cannot read enum %s", name)
		}
		if len(got) != len(labels) {
			return fmt.Errorf("incompatible enum %s; export labels before deployment", name)
		}
		for _, label := range labels {
			ok := false
			for _, v := range got {
				if v == label {
					ok = true
				}
			}
			if !ok {
				return fmt.Errorf("incompatible enum %s", name)
			}
		}
	}
	return nil
}
