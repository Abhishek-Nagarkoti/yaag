/*
 * This is the main core of the yaag package
 */
package yaag

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

const TEMPLATE = `<!DOCTYPE html>
<html>
<head lang="en">
    <title> API Documentation </title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
    <script src="http://google-code-prettify.googlecode.com/svn/loader/run_prettify.js"></script>
    <link href='http://fonts.googleapis.com/css?family=Roboto' rel='stylesheet' type='text/css'>
    <link rel="stylesheet" href="http://google-code-prettify.googlecode.com/svn/trunk/src/prettify.css">
    <!-- Optional theme -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <!-- Latest compiled and minified JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
    <style type="text/css">
        body {
            font-family: 'Roboto', sans-serif;
        }
    </style>
    <style type="text/css">
        pre.prettyprint {
            border: 1px solid #ccc;
            margin-bottom: 0;
            padding: 9.5px;
        }
    </style>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
    <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse"
                    data-target="#bs-example-navbar-collapse-1">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="#">{{.Title}}</a>
        </div>

        <!-- Collect the nav links, forms, and other content for toggling -->
        <div class="collapse navbar-collapse pull-right" id="bs-example-navbar-collapse-1">
            <form class="navbar-form navbar-left" role="search">
                <div class="form-group">
                    <input type="text" class="form-control" placeholder="Search">
                </div>
                <button type="submit" class="btn btn-default">Find</button>
            </form>
        </div>
        <!-- /.navbar-collapse -->
    </div>
    <!-- /.container-fluid -->
</nav>
<div class="container" style="margin-top: 70px;">
    <div class="alert alert-info">
        <p>Base URL => <strong>{{.BaseLink}}</strong></p></div>
    <hr>
    {{ range $key, $value := .array }}
    <h4 id="{{$key}}top"><a class="anchor" href="#{{$key}}top"><span class="glyphicon glyphicon-link"
              aria-hidden="true"></span></a> <code>{{$value.HttpVerb}}
        {{$value.Path}}</code></h4>
    {{ range $wrapperKey, $wrapperValue := $value.HtmlValues }}
    <div id="{{mult $key $wrapperValue.Id}}next" class="container" style="margin-left:2em;">
        <h4  style="cursor:pointer;" type="button" data-toggle="collapse" data-target="#{{mult $key $wrapperValue.Id}}container"
            aria-expanded="false" aria-controls="collapseExample"><a class="anchor" href="#{{mult $key $wrapperValue.Id}}next"><span class="glyphicon glyphicon-link"
                                                                        aria-hidden="true"></span></a> Example {{add $wrapperKey 1}}
        </h4>

        <div class="collapse" id="{{mult $key $wrapperValue.Id}}container">
            {{ if $wrapperValue.RequestHeader }}
            <p> <H4> Request Headers </H4> </p>
            <table class="table table-bordered table-striped">
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
                {{ range $key, $value := $wrapperValue.RequestHeader }}
                <tr>
                    <td>{{ $key }}</td>
                    <td> {{ $value }}</td>
                </tr>
                {{ end }}
            </table>
            {{ end }}

            {{ if $wrapperValue.PostForm }}
            <p> <H4> Post Form </H4> </p>
            <table class="table table-bordered table-striped">
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
                {{ range $key, $value := $wrapperValue.PostForm }}
                <tr>
                    <td>{{ $key }}</td>
                    <td> {{ $value }}</td>
                </tr>
                {{ end }}
            </table>
            {{ end }}


            {{ if $wrapperValue.RequestUrlParams }}
            <p> <H4> URL Params </H4> </p>
            <table class="table table-bordered table-striped">
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
                {{ range $key, $value := $wrapperValue.RequestUrlParams }}
                <tr>
                    <td>{{ $key }}</td>
                    <td> {{ $value }}</td>
                </tr>
                {{ end }}
            </table>
            {{ end }}

            {{ if $wrapperValue.RequestBody }}
            <p> <H4> Request Body </H4> </p>
            <pre class="prettyprint lang-json">{{ $wrapperValue.RequestBody }}</pre>
            {{ end }}

            <p><h4> Response Code</h4></p>
            <pre class="prettyprint lang-json">{{ $wrapperValue.ResponseCode }}</pre>

            {{ if $wrapperValue.ResponseHeader }}
            <p><h4> Response Headers</h4></p>
            <table class="table table-bordered table-striped">
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
                {{ range $key, $value := $wrapperValue.ResponseHeader }}
                <tr>
                    <td>{{ $key }}</td>
                    <td> {{ $value }}</td>
                </tr>
                {{ end }}
            </table>
            {{ end }}


            {{ if $wrapperValue.ResponseBody }}
            <p> <H4> Response Body </H4> </p>
            <pre class="prettyprint lang-json">{{ $wrapperValue.ResponseBody }}</pre>
            {{ end }}
        </div>
    </div>
    <hr>
    {{ end}}
    {{ end}}
</div>
</body>
</html>`

var CommonHeaders = []string{
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Cache-Control",
	"Content-Length",
	"Content-Type",
	"Origin",
	"User-Agent",
	"X-Forwarded-For",
}

type APICall struct {
	Id int

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestUrlParams     map[string]string

	RequestBody  string
	ResponseBody string
	ResponseCode int
}

type PathSpec struct {
	HttpVerb   string
	Path       string
	HtmlValues []APICall
}

type ApiCallValue struct {
	BaseLink string
	Path     []PathSpec
}

type Config struct {
	Init     bool
	DocTitle string
	DocPath  string
}

var ApiCallValueInstance = &ApiCallValue{}

func add(x, y int) int {
	return x + y
}

func mult(x, y int) int {
	return x * y
}

func main() {
	//	firstApi := APICall{Id: 1, MethodType: "GET", CurrentPath: "/login/:id", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},

	//		RequestBody: "{ 'main' : { 'id' : 2, 'name' : 'Gopher' }}"}

	secondApi := APICall{Id: 2, MethodType: "POST", CurrentPath: "/singup", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		ResponseBody: "{ 'main' : { 'Key' : 'ABC-123-XYZ', 'name' : 'Gopher' }}"}

	config := Config{Init: false, DocPath: "html/home.html", DocTitle: "YAAG"}

	//	valueArray := []APICall{secondApi, firstApi}
	//	allApis := ApiCallValue{BaseLink: "www.google.com", HtmlValues: valueArray}
	GenerateHtml(&secondApi, &config)
}

func GenerateHtml(htmlValue *APICall, config *Config) {
	shouldAddPathSpec := true
	log.Printf("PathSpec : %v", ApiCallValueInstance.Path)
	for _, pathSpec := range ApiCallValueInstance.Path {
		if pathSpec.Path == htmlValue.CurrentPath && pathSpec.HttpVerb == htmlValue.MethodType {
			shouldAddPathSpec = false
			shouldAdd := true
			for _, value := range pathSpec.HtmlValues {
				if value.ResponseBody == htmlValue.ResponseBody {
					shouldAdd = false
				}
			}
			if shouldAdd {
				htmlValue.Id = len(pathSpec.HtmlValues) + 1
				pathSpec.HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
			}
		}
	}

	if shouldAddPathSpec {
		pathSpec := PathSpec{
			HttpVerb: htmlValue.MethodType,
			Path:     htmlValue.CurrentPath,
		}
		htmlValue.Id = 1
		pathSpec.HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
		ApiCallValueInstance.Path = append(ApiCallValueInstance.Path, pathSpec)
	}
	funcs := template.FuncMap{"add": add, "mult": mult}
	t := template.New("API Documentation").Funcs(funcs)
	filePath, err := filepath.Abs(config.DocPath)
	htmlString := TEMPLATE

	t, err = t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()

	if err != nil {
		log.Println(err)
		return
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{"array": ApiCallValueInstance.Path,
		"BaseLink": ApiCallValueInstance.BaseLink, "Title": config.DocTitle})
}
