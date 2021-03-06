package archive

import (
	"fmt"
	"os"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccArchiveFile_Basic(t *testing.T) {
	var fileSize string
	r.Test(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccArchiveFileContentConfig,
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists("zip_file_acc_test.zip", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			r.TestStep{
				Config: testAccArchiveFileFileConfig,
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists("zip_file_acc_test.zip", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			r.TestStep{
				Config: testAccArchiveFileDirConfig,
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists("zip_file_acc_test.zip", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			r.TestStep{
				Config: testAccArchiveFileOutputPath,
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(fmt.Sprintf("%s/test.zip", tmpDir), &fileSize),
				),
			},
		},
	})
}

func testAccArchiveFileExists(filename string, fileSize *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		*fileSize = ""
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		*fileSize = fmt.Sprintf("%d", fi.Size())
		return nil
	}
}

var testAccArchiveFileContentConfig = `
data "archive_file" "foo" {
  type                    = "zip"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "zip_file_acc_test.zip"
}
`

var tmpDir = os.TempDir() + "/test"
var testAccArchiveFileOutputPath = fmt.Sprintf(`
data "archive_file" "foo" {
  type                    = "zip"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "%s/test.zip"
}
`, tmpDir)

var testAccArchiveFileFileConfig = `
data "archive_file" "foo" {
  type        = "zip"
  source_file = "test-fixtures/test-file.txt"
  output_path = "zip_file_acc_test.zip"
}
`

var testAccArchiveFileDirConfig = `
data "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir"
  output_path = "zip_file_acc_test.zip"
}
`
