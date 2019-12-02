version: '2'

services:{{range $val := .Index}}
    node{{$val}}:
        image: docker.io/centos:7.5.1804
        command: ["bin/start_xchain.sh"]
        working_dir: /var/node{{$val}}
        volumes:
            - {{$.SandRoot}}/nodes/node{{$val}}:/var/node{{$val}}
            - {{$.SandRoot}}/bin:/var/node{{$val}}/bin
        ports:{{if lt $val $.PortSeg}}
            - 3710{{$val}}:37101
            - 4710{{$val}}:47101{{else if ge $val $.PortSeg}}
            - 371{{$val}}:37101
            - 471{{$val}}:47101{{end}}
{{end}}
