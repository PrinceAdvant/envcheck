package envcheck

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	tpl, _ := Parse(strings.NewReader("# db\nDATABASE_URL=\nREDIS_URL=\nexport MAIL_FROM=\"\"\n"))
	env, _ := Parse(strings.NewReader("DATABASE_URL=postgres://x\nREDIS_URL=\nDEBUG='1'\n"))
	r := Compare(tpl, env)
	if strings.Join(r.Missing, ",") != "MAIL_FROM" || strings.Join(r.Empty, ",") != "REDIS_URL" || strings.Join(r.Extra, ",") != "DEBUG" {
		t.Fatalf("unexpected report %+v", r)
	}
	if r.OK() {
		t.Fatal("expected not OK")
	}
	if env["DEBUG"] != "1" {
		t.Fatal("quotes not stripped")
	}
}
