FROM debian:bookworm-slim as setup

RUN apt update && apt-get install -y wget unzip

RUN wget https://powerwalker.com/wp-content/uploads/2022/01/221019-PowerWalker_setup_LinuxAMD64.tar.zip
RUN unzip /221019-PowerWalker_setup_LinuxAMD64.tar.zip
RUN tar -xvzf "/221019 PowerWalker_setup_LinuxAMD64.tar.gz"
RUN /PowerWalker_setup_LinuxAMD64/LinuxAMD64/install.bin -i silent; exit 0

FROM debian:bookworm-slim

COPY --from=setup /opt/MonitorSoftware /opt/MonitorSoftware
COPY config.properties /opt/MonitorSoftware/config/config.properties
COPY run.sh /run.sh
RUN chmod +x /run.sh

ENTRYPOINT ["/run.sh"]