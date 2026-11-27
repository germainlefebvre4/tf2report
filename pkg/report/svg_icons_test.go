package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testdataIcons(name string) string {
	return filepath.Join("..", "..", "testdata", "icons", name)
}

func TestBuiltinGlyphs(t *testing.T) {
	pack, err := loadBuiltin()
	if err != nil {
		t.Fatalf("loadBuiltin: %v", err)
	}
	for _, id := range []string{"compute", "network", "storage", "database", "security", "generic"} {
		ic, ok := pack.icons[id]
		if !ok {
			t.Errorf("built-in glyph %s is missing", id)
			continue
		}
		if ic.viewBox != "0 0 24 24" || ic.inner == "" {
			t.Errorf("glyph %s = %+v", id, ic)
		}
	}
	for _, id := range []string{"aws", "azure", "gcp", "aws-ec2"} {
		if _, ok := pack.icons[id]; ok {
			t.Errorf("no built-in glyph may exist for %s", id)
		}
	}
}

func TestLoadPack_Valid(t *testing.T) {
	var skipped []string
	pack, err := loadPack(testdataIcons("valid"), func(f, _ string) { skipped = append(skipped, f) })
	if err != nil {
		t.Fatalf("loadPack: %v", err)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v, want none", skipped)
	}
	for _, id := range []string{"aws-ec2", "aws", "azure-vm", "compute", "generic"} {
		if _, ok := pack.icons[id]; !ok {
			t.Errorf("pack misses %s", id)
		}
	}
	if _, ok := pack.icons["aws-vpc"]; ok {
		t.Errorf("nested directory must be ignored")
	}
	if len(pack.icons) != 5 {
		t.Errorf("pack has %d icons, want 5", len(pack.icons))
	}
	if g := pack.icons["generic"].inner; !strings.Contains(g, `id="generic-g"`) || !strings.Contains(g, "url(#generic-g)") {
		t.Errorf("ids are not prefixed by the family id: %s", g)
	}
}

func TestLoadPack_Hostile(t *testing.T) {
	skipped := map[string]string{}
	pack, err := loadPack(testdataIcons("hostile"), func(f, reason string) { skipped[f] = reason })
	if err != nil {
		t.Fatalf("loadPack: %v", err)
	}
	if len(pack.icons) != 0 {
		t.Errorf("hostile pack accepted %d icons", len(pack.icons))
	}
	entries, _ := os.ReadDir(testdataIcons("hostile"))
	for _, e := range entries {
		if skipped[e.Name()] == "" {
			t.Errorf("%s was not rejected with a reason", e.Name())
		}
	}
}

func TestLoadPack_Oversize(t *testing.T) {
	dir := t.TempDir()
	big := `<svg xmlns="http://www.w3.org/2000/svg"><title>` + strings.Repeat("a", maxIconSize) + `</title></svg>`
	if err := os.WriteFile(filepath.Join(dir, "big.svg"), []byte(big), 0o644); err != nil {
		t.Fatal(err)
	}
	var reason string
	pack, err := loadPack(dir, func(_, r string) { reason = r })
	if err != nil || len(pack.icons) != 0 || !strings.Contains(reason, "larger") {
		t.Errorf("oversize file: err=%v icons=%d reason=%q", err, len(pack.icons), reason)
	}
}

func TestLoadPack_Errors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "afile")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(t.TempDir(), "missing")
	for _, dir := range []string{missing, file} {
		t.Run(dir, func(t *testing.T) {
			_, err := loadPack(dir, nil)
			if err == nil || !strings.Contains(err.Error(), dir) {
				t.Errorf("error = %v, want one naming %s", err, dir)
			}
		})
	}

	t.Run("unreadable", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "locked")
		if err := os.Mkdir(dir, 0o000); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		if os.Geteuid() == 0 {
			t.Skip("root can read any directory")
		}
		if _, err := loadPack(dir, nil); err == nil || !strings.Contains(err.Error(), dir) {
			t.Errorf("error = %v, want one naming %s", err, dir)
		}
	})

	t.Run("no usable file", func(t *testing.T) {
		pack, err := loadPack(t.TempDir(), nil)
		if err != nil || len(pack.icons) != 0 {
			t.Errorf("empty dir: err=%v icons=%d", err, len(pack.icons))
		}
	})
}

func TestIconSet_Resolve(t *testing.T) {
	set, err := newIconSet(testdataIcons("valid"), nil)
	if err != nil {
		t.Fatalf("newIconSet: %v", err)
	}
	tests := []struct {
		name  string
		chain []string
		want  string
		inPak bool
	}{
		{"service in pack", []string{"aws-ec2", "aws", "compute", "generic"}, "aws-ec2", true},
		{"provider level in pack", []string{"aws-foo", "aws", "storage", "generic"}, "aws", true},
		{"pack wins over built-in", []string{"compute", "generic"}, "compute", true},
		{"built-in category", []string{"storage", "generic"}, "storage", false},
		{"generic from pack", []string{"zzz", "generic"}, "generic", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ic, ok := set.resolve(tt.chain)
			if !ok || id != tt.want {
				t.Fatalf("resolve = %q, %v, want %q", id, ok, tt.want)
			}
			if _, fromPack := set.pack.icons[id]; fromPack != tt.inPak || ic.inner == "" {
				t.Errorf("fromPack = %v, want %v", fromPack, tt.inPak)
			}
		})
	}

	builtinOnly, err := newIconSet("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if id, _, _ := builtinOnly.resolve([]string{"aws-foo", "aws", "storage", "generic"}); id != "storage" {
		t.Errorf("without pack aws falls to category glyph, got %q", id)
	}
}

func TestSanitizeIcon_Reserialized(t *testing.T) {
	src := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><!-- c --><rect  x="1"   y="2" width="3" height="4"/></svg>`
	ic, err := sanitizeIcon([]byte(src), "fam")
	if err != nil {
		t.Fatal(err)
	}
	if ic.viewBox != "0 0 10 10" || ic.inner != `<rect x="1" y="2" width="3" height="4"></rect>` {
		t.Errorf("icon = %+v", ic)
	}
}

func TestVendorIcons(t *testing.T) {
	vendor, err := loadVendor()
	if err != nil {
		t.Fatalf("loadVendor: %v", err)
	}
	counts := map[string]int{}
	for id := range vendor.icons {
		counts[id[:strings.Index(id, "-")]]++
	}
	for cloud, min := range map[string]int{"aws": 20, "azure": 20, "gcp": 5} {
		if counts[cloud] < min {
			t.Errorf("%s vendor icons = %d, want at least %d", cloud, counts[cloud], min)
		}
	}
}

func TestIconSet_VendorBetweenPackAndGlyph(t *testing.T) {
	chain := []string{"aws-ec2", "aws", "compute", "generic"}

	plain, err := newIconSet("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if id, _, _ := plain.resolve(chain); id != "aws-ec2" {
		t.Errorf("without pack resolved %q, want the vendor icon aws-ec2", id)
	}
	if id, _, _ := plain.resolve([]string{"aws-unknown", "aws", "compute", "generic"}); id != "compute" {
		t.Errorf("unknown service resolved %q, want the compute glyph", id)
	}

	withPack, err := newIconSet(testdataIcons("valid"), nil)
	if err != nil {
		t.Fatal(err)
	}
	id, ic, _ := withPack.resolve(chain)
	if _, fromPack := withPack.pack.icons[id]; !fromPack || ic.inner != withPack.pack.icons["aws-ec2"].inner {
		t.Errorf("user pack must win over the vendor icon")
	}
}
