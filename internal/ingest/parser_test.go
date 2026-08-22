package ingest

import (
	"strings"
	"testing"
)

func TestParseHTMLRemovesNoiseAndExtractsMetadata(t *testing.T) {
	raw := `<!doctype html><html><head><title>Aryakun Services</title><meta name="description" content="Web development services"><link rel="canonical" href="/services"></head><body><header>Menu</header><nav>Home Services Contact</nav><main><h1>Website Development</h1><p>Kami membuat website modern untuk bisnis.</p><script>alert('x')</script><p>Hubungi kami untuk informasi lebih lanjut.</p></main><footer>Copyright</footer></body></html>`
	page, err := ParseHTML(raw, "https://aryakun.id/services")
	if err != nil { t.Fatalf("ParseHTML() error = %v", err) }
	if page.Title != "Aryakun Services" { t.Fatalf("Title = %q", page.Title) }
	if page.Description != "Web development services" { t.Fatalf("Description = %q", page.Description) }
	if page.Canonical != "https://aryakun.id/services" { t.Fatalf("Canonical = %q", page.Canonical) }
	if !strings.Contains(page.Content, "Website Development") || !strings.Contains(page.Content, "website modern") { t.Fatalf("content missing readable text: %q", page.Content) }
	for _, noise := range []string{"alert", "Copyright", "Home Services Contact", "Menu"} { if strings.Contains(page.Content, noise) { t.Errorf("content contains noise %q: %q", noise, page.Content) } }
}

func TestParseHTMLRejectsInsufficientContent(t *testing.T) {
	_, err := ParseHTML("<html><body><h1>Hi</h1></body></html>", "https://aryakun.id/")
	if err == nil { t.Fatal("ParseHTML() expected insufficient-content error") }
}

func TestParsedPageDocument(t *testing.T) {
	page := ParsedPage{Title:"Test", Canonical:"https://aryakun.id/test", Content:"A sufficiently long page content for indexing."}
	doc := page.Document("doc-1", "https://aryakun.id/test")
	if doc.Source != "website" || doc.Type != "webpage" || doc.Trust != "first_party" { t.Fatalf("unexpected document metadata: %+v", doc) }
	if doc.ContentHash == "" { t.Fatal("document hash is empty") }
}
