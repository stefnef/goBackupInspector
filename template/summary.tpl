{{template "header"}}
{{$LeftDir := .LeftDir}}
{{$RightDir := .RightDir}}
{{ $lengthFilesNotInLeftDir := ShowFilesNotInLeftDir . | len}}
{{ $lengthFilesNotInRightDir := ShowFilesNotInRightDir . | len}}
{{ $lengthDirsNotInLeftDir := ShowDirsNotInLeftDir . | len}}
{{ $lengthDirsNotInRightDir := ShowDirsNotInRightDir . | len}}
{{ $lengthUnequalFiles := len .UnequalFiles }}
{{ $lengthIgnoredElements := len .IgnoredElement }}
{{ $lengthComparedFiles := len .ComparedFiles }}
<h1>Summary Of Comparison</h1>
<body>
<div class="container-fluid">
    <div class="row">
        <ul class="list-group  collapse in"  id="MissingFilesLeft">
            <li class="list-group-item">Backup File: {{.BackupFileName}}</li>
            <li class="list-group-item">Creation Date: {{ .Date.Format "02.01.2006"}}</li>
            <li class="list-group-item">With Differences: {{.WithDifferences}}</li>
        </ul>
    </div>
    <div class="row">
        <h2>Missing Files/Directories</h2>
        <div class="col-sm-6">
            <div class="container-fluid">
                <h3>{{ .LeftDir }}</h3>
                <div class="container-fluid" data-toggle="collapse" data-target="#MissingFilesLeft">
                    {{ if gt $lengthFilesNotInLeftDir 0 }}
                        <h4>These files are missing in {{ .LeftDir }}</h4>
                        <ul class="list-group  collapse in"  id="MissingFilesLeft">
                        {{ range ShowFilesNotInLeftDir .}}
                            <li class="list-group-item">{{ . }}</li>
                        {{end}}
                        </ul>
                    {{else}}
                        <h4>No files are missing in {{ .LeftDir }}</h4>
                    {{end}}
                </div>

                <div class="container-fluid" data-toggle="collapse" data-target="#MissingDirsLeft">
                    {{ if gt $lengthDirsNotInLeftDir 0 }}
                        <h4>These directories are missing in {{ .LeftDir }}</h4>
                        <ul class="list-group  collapse in"  id="MissingDirsLeft">
                            {{ range ShowDirsNotInLeftDir .}}
                                <li class="list-group-item">{{ . }}</li>
                            {{end}}
                        </ul>
                    {{else}}
                        <h4>No directories are missing in {{ .LeftDir }}</h4>
                    {{end}}
                </div>
            </div>
        </div>
        <div class="col-sm-6">
            <div class="container-fluid">
                <h3>{{ .RightDir }}</h3>
                <div class="container-fluid" data-toggle="collapse" data-target="#MissingFilesRight">
                    {{ if gt $lengthFilesNotInRightDir 0 }}
                        <h4>These files are missing in {{ .RightDir }}</h4>
                        <ul class="list-group  collapse in"  id="MissingFilesRight">
                            {{ range ShowFilesNotInRightDir .}}
                                <li class="list-group-item">{{ . }}</li>
                            {{end}}
                        </ul>
                    {{else}}
                        <h4>No files are missing in {{ .RightDir }}</h4>
                    {{end}}
                </div>

                <div class="container-fluid" data-toggle="collapse" data-target="#MissingDirsRight">
                    {{ if gt $lengthDirsNotInRightDir 0 }}
                        <h4>These directories are missing in {{ .RightDir }}</h4>
                        <ul class="list-group  collapse in"  id="MissingDirsRight">
                            {{ range ShowDirsNotInRightDir .}}
                                <li class="list-group-item">{{ . }}</li>
                            {{end}}
                        </ul>
                    {{else}}
                        <h4>No directories are missing in {{ .RightDir }}</h4>
                    {{end}}
                </div>
            </div>
        </div>
    </div>
    <hr>
    <div class="row">
        <div class="container-fluid">
            {{ if gt $lengthUnequalFiles 0 }}
                <h2 data-toggle="collapse" data-target="#UnequalFiles">Unequal files:</h2>
                <ul class="list-group  collapse in"  id="UnequalFiles">
                    {{ range .UnequalFiles }}
                        <li class="list-group-item">{{ . }}</li>
                    {{end}}
                </ul>
            {{else}}
                <h2>No Unequal Files Found!</h2>
            {{end}}
        </div>
    </div>
    <hr>
    <div class="row">
        <div class="container-fluid">
            {{ if gt $lengthIgnoredElements 0 }}
            <h2 data-toggle="collapse" data-target="#IgnoredElements">Ignored elements:</h2>
            <ul class="list-group  collapse in"  id="IgnoredElements">
                {{ range .IgnoredElement }}
                    <li class="list-group-item">{{ . }}</li>
                {{end}}
            </ul>
            {{else}}
                <h2>No Ignored Elements</h2>
            {{end}}
        </div>
    </div>
    <hr>
    <div class="row">
        <div class="container-fluid">
            {{ if gt $lengthComparedFiles 0 }}
                <h2 data-toggle="collapse" data-target="#ComparedFiles">Compared Files:</h2>
                <ul class="list-group  collapse"  id="ComparedFiles">
                    {{ range .ComparedFiles }}
                        <li class="list-group-item">{{ . }}</li>
                    {{end}}
                </ul>
            {{else}}
                <h2>No Files Were Compared!</h2>
                    <p>Hm, this sounds strange...</p>
            {{end}}
        </div>
    </div>
</div>
{{template "footer"}}