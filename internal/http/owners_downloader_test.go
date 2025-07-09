package http

import (
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/mrclmr/icm/internal/cont"
)

func Test_parseOwners(t *testing.T) {
	tests := []struct {
		name    string
		body    io.Reader
		want    []cont.Owner
		wantErr bool
	}{
		{
			"Parsing valid HTML body returns owners map",
			validBody(),
			[]cont.Owner{
				{
					Code:    "AAA",
					Company: "A Company",
					City:    "A City",
					Country: "A Country",
				},
				{
					Code:    "BBB",
					Company: "B Company",
					City:    "B City",
					Country: "B Country",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOwners(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOwners() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("parseOwners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func validBody() io.Reader {
	return strings.NewReader(`<!DOCTYPE html>
<body>
	<main id="search" role="main">

        <section id="archive" class="searchResponse">
            <div class="container-fluid g-0">
                <div class="row">
                    <div class="col-sm-12">
                        <div class="container-xxxl g-3 g-lg-5">
                            <div class="row">
                                <section class="breadCrumbSection mb-3 relative noBg d-flex justify-content-end">
                                    <div class="col-sm-12 mt-2 hideOnMobile">
                                        <ul id="breadcrumbs" class="breadcrumbs">
                                            <li class="item-home">
                                                <a class="bread-link bread-home" href="https://www.bic-code.org"
                                                   title="Home">Home</a>
                                            </li>
                                            <li class="separator separator-home">
                                                |
                                            </li>
                                            <li class="item-home">
                                                <a class="bread-link bread-home" href="/bic-codes" title="Bic Codes">Bic
                                                    Codes</a>
                                            </li>
                                            <li class="separator separator-home">
                                                |
                                            </li>
                                            <li class="item-current">
                                                <strong class="bread-current" title="Search Results">Search
                                                    Results</strong>
                                            </li>
                                        </ul>
                                    </div>

                                    <a class="btn-primary searchBtn absolute" href="/bic-codes">SEARCH AGAIN</a>
                                </section>
                                <h1 class="text-secondary">BIC Code Search Results</h1>
                                <div class="d-flex align-items-center justify-content-between">
                                    <p class="resultsText">3675 search results for "<span class="upperCase">all"</span>
                                    </p>
                                </div>

                                <div class="col-sm-12">
                                    <div id="bicloader" class="justify-content-center">
                                        <img class="no-lazy"
                                             src="https://www.bic-code.org/wp-content/themes/tessellate/assets/images/bicloader.gif"
                                             alt>
                                    </div>
                                    <table width="100%" id="bic-datatable" class="bicResults mt-1 responsive"
                                           data-type="bic">
                                        <thead>
                                        <tr>
                                            <th class="desktop tablet-l tablet-p mobile-l mobile-p">Code</th>
                                            <th class="min-tablet-l">Company</th>
                                            <th class="min-tablet-l" name="address">Address</th>
                                            <th class="min-desktop">City</th>
                                            <th class="min-desktop">Zip</th>
                                            <th class="min-tablet-l">Country</th>
                                            <th class="desktop tablet-l">Details</th>
                                        </tr>
                                        </thead>
                                        <tbody>
                                        <tr class="nostripe">
                                            <td class="flexMobile align-items-center dtr-control sorting_1" tabindex="0"></span></td>
                                            <td class="flexMobile align-items-center">VANTAGE FORWARDING &amp; SHIPPING HK LIMITED</td>
                                            <td class="flexMobile align-items-center">33B United Centre, 95 Queensway, Admiralty</td>
                                            <td class="flexMobile align-items-center"></td>
                                            <td class="flexMobile align-items-center">200072</td>
                                            <td class="flexMobile align-items-center">HK</td>
                                            <td data-geofence="0" class="no-sort flexMobile detailWidth">
                                              <div class="d-flex my-1 my-md-0 position-relative d-flex justify-content-start justify-content-md-end align-items-center">
                                                <a class="upperCase withArrow" href="/bic-codes//">View</a>
                                              </div>
                                            </td>
                                        </tr>
                                        <tr class="nostripe">
                                            <td class="flexMobile align-items-center">AAAU</span></td>
                                            <td class="flexMobile align-items-center">A Company</span></td>
                                            <td class="flexMobile align-items-center"></span></td>
                                            <td class="flexMobile align-items-center">A City</span></td>
                                            <td class="flexMobile align-items-center"></span></td>
                                            <td class="flexMobile align-items-center">A Country</span></td>
                                            <td class="no-sort flexMobile detailWidth align-items-center">
                                                <a class="upperCase withArrow" href="/bic-codes/aaau/">View</a>
                                            </td>
                                        </tr>
<tr class="nostripe">
                                            <td class="flexMobile align-items-center">BBBU</span></td>
                                            <td class="flexMobile align-items-center">B Company</span></td>
                                            <td class="flexMobile align-items-center"></span></td>
                                            <td class="flexMobile align-items-center">B City</span></td>
                                            <td class="flexMobile align-items-center"></span></td>
                                            <td class="flexMobile align-items-center">B Country</span></td>
                                            <td class="no-sort flexMobile detailWidth align-items-center">
                                                <a class="upperCase withArrow" href="/bic-codes/bbbu/">View</a>
                                            </td>
                                        </tr>
                                        </tbody>
                                        <tfoot></tfoot>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

    </main>
</body>
</html>
`)
}
