version: '2'

services:
    node1:
        image: docker.io/centos:7.5.1804
        command: ["bin/start_xchain.sh"]
        working_dir: {{root}}/node1
        volumes: 
            - {{root}}/node1:{{root}}/node1
            - {{root}}/bin:{{root}}/node1/bin
            - /var/run/docker.sock:/var/run/docker.sock
        ports: 
            - 27101:37101
        user: "{{user}}"

    node2:
        image: docker.io/centos:7.5.1804
        command: ["bin/start_xchain.sh"]
        working_dir: {{root}}/node2
        volumes: 
            - {{root}}/node2:{{root}}/node2
            - {{root}}/bin:{{root}}/node2/bin
            - /var/run/docker.sock:/var/run/docker.sock
        ports: 
            - 27102:37101
        user: "{{user}}"

    node3:
        image: docker.io/centos:7.5.1804
        command: ["bin/start_xchain.sh"]
        working_dir: {{root}}/node3
        volumes: 
            - {{root}}/node3:{{root}}/node3
            - {{root}}/bin:{{root}}/node3/bin
            - /var/run/docker.sock:/var/run/docker.sock
        ports: 
            - 27103:37101
        user: "{{user}}"

    cli:
        image: docker.io/centos:7.5.1804
        command: "bash"
        working_dir: {{root}}/client
        volumes: 
            - {{root}}/bin:{{root}}/bin
            - {{root}}/client:{{root}}/client
        stdin_open: true
        tty: true
        user: "{{user}}"
        

