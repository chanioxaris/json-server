package common

import (
	"html/template"
	"net/http"

	"github.com/chanioxaris/json-server/storage"
	"github.com/chanioxaris/json-server/web"
)

const homePageTemplate = `
<!doctype html>
<html lang="en">
	<head>
    	<!-- Required meta tags -->
    	<meta charset="utf-8">
    	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    	<!-- Bootstrap CSS -->
    	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

    	<title>JSON Server</title>

		<style> 
			.tooltip-inner {
				text-align: left;
			}
		</style>
  	</head>
  	<body>
		<div class="container mt-5">
			<h1>Welcome to JSON Server</h1>

			<p>Below you can find a list of the available resources generated from the provided json file.</p>
			<p>You can hover over on any resource to see the available endpoints.</p>
			</br>

			<h2>Resources</h2>
			{{ range $resource, $singular := . }}
				{{ if $singular }}
					/{{ $resource }}
					<span 
						class="badge badge-secondary"
						data-toggle="tooltip" 
						data-html="true"
						data-placement="right" 
						title="<ul><li>GET /{{ $resource }}</li></ul>"
					>
						1
					</span>
				{{ else }}
					/{{ $resource }}
					<span 
						class="badge badge-secondary"
						data-toggle="tooltip" 
						data-html="true"
						data-placement="right" 
						title="<ul><li>GET /{{ $resource }}</li><li>GET /{{ $resource }}/:id</li><li>POST /{{ $resource }}</li><li>PUT /{{ $resource }}/:id</li><li>PATCH /{{ $resource }}/:id</li><li>DELETE /{{ $resource }}/:id</li></ul>"
					>
						5
					</span>
				{{ end }}
				</br>
			{{ end }}
		</div>

		<footer class="fixed-bottom text-center mb-3">
			<a href="https://github.com/chanioxaris/json-server" target="_blank">
				github.com/chanioxaris/json-server
			</a>
		</footer>

    	<!-- Optional JavaScript -->
    	<!-- jQuery first, then Popper.js, then Bootstrap JS -->
    	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
    	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
    	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  	
		<!-- Enable tooltips -->
		<script type="text/javascript">
			$(document).ready(function () {
			  	$('[data-toggle="tooltip"]').tooltip()
			});
		</script>
	</body>
</html>
`

// HomePage renders the home page template with useful information about generated endpoints and resources.
func HomePage(storageResources map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("home").Parse(homePageTemplate)
		if err != nil {
			web.Error(w, http.StatusBadRequest, storage.ErrInternalServerError.Error())
			return
		}

		if err = t.Execute(w, storageResources); err != nil {
			web.Error(w, http.StatusBadRequest, storage.ErrInternalServerError.Error())
			return
		}
	}
}
