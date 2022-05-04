FROM rundeck/rundeck:3.4.3

# update config
RUN sed -i 's|csp/enabled", "true"|csp/enabled", "false"|g' /etc/remco/templates/rundeck-config-csp.properties
RUN sed -i 's|duration/max", "30d"|duration/max", "0"|g' /etc/remco/templates/rundeck-config.properties
