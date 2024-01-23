#!/bin/bash
cd /opt/MonitorSoftware && ./agent start
tail -f /opt/MonitorSoftware/mymanager.log