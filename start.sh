#!/usr/bin/env bash

source ./config.sh
source ./common/colors.sh
source ./common/spinner.sh
source ./common/email.sh
source ./common/secret.sh
source ./common/display_header.sh
source ./runner.sh

display_header
echo ""

echo -e "$(set_color "$(set_color " start " "$BOLD_INTNS_WHITE")" "$INTNS_BG_BLUE") Starting pipeline"

NOTIFICATION_TITLE=""
NOTIFICATION_BODY=""
SECONDS=0

chmod u+x ./tasks/*

for step in tasks/*.sh; do
    run_with_spinner "$(basename $step)" "$step"
    runner_result=$?

    if [ $runner_result -ne 0 ]; then
        printf "\a"
        printf "\a"
        printf "\a"

        duration=$SECONDS
        CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

        NOTIFICATION_TITLE="<h2 style='color:#e81a1a; font-size:26px; margin:0 0 15px 0;'>Workflow failed! &#127755;</h2><p style='font-size:17px; line-height:1.5; color:#c8c8c888;'>Pipeline ran for $((duration / 60)) minutes and $((duration % 60)) seconds.<p/><p style='color:#9a9a9a;margin:0 0 25px 0;'>Logs of this workflow can be found in $CURRENT_WORKFLOW_LOG_NAME.txt</p>"
        NOTIFICATION_BODY="${NOTIFICATION_BODY}<tr><td colspan='2' style='width: 100%; border: 1px solid #474747'></td></tr><tr><td><p style='margin:0; color:#e81a1a;'><strong>FAILED</strong></p></td><td><p>&#10060;</p></td></tr>"
        
        echo "[$(date --utc +%FT%TZ)]   WORKFLOW FINISHED WITH AN ERROR!" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
        echo -e "$(set_color "$(set_color "  stop " "$BOLD_INTNS_WHITE")" "$INTNS_BG_RED") Pipeline failed and ran for $((duration / 60)) minutes and $((duration % 60)) seconds."
        
        echo ""
        echo -e $(set_color "Logs are available at $CURRENT_WORKFLOW_LOG_NAME" "$BOLD_INTNS_RED")

        if [ $NOTIFICATIONS_ENABLED -eq 1 ]; then 
            EMAIL=$(create_workflow_notification_email "$NOTIFICATION_TITLE" "$NOTIFICATION_BODY")
            send_email "$EMAIL_TO" "Bash CI/CD fail" "$EMAIL" "$(get_secret key resend_key)"
        fi

        exit 1
    else
        printf "\a"
    fi
done

chmod u-x ./tasks/*

duration=$SECONDS
CURRENT_WORKFLOW_LOG_NAME=$(ls -t "$TEMP_DIR" | head -1)

NOTIFICATION_TITLE="<h2 style='color:#1a73e8; font-size:26px; margin:0 0 15px 0;'>Workflow completed! &#128507;</h2><p style='font-size:17px; line-height:1.5; color:#c8c8c888;'>Pipeline ran for $((duration / 60)) minutes and $((duration % 60)) seconds.</p><p style='color:#9a9a9a;margin:0 0 25px 0;'>Logs of this workflow can be found in $CURRENT_WORKFLOW_LOG_NAME.txt</p>"
NOTIFICATION_BODY="${NOTIFICATION_BODY}<tr><td colspan='2' style='width: 100%; border: 1px solid #474747'></td></tr><tr><td><p style='margin:0; color:#1a73e8;'><strong>COMPLETED</strong></p></td><td><p>&#9989;</p></td></tr>"

echo "[$(date --utc +%FT%TZ)]   WORKFLOW COMPLETED!" &>> "$LOGS_DIR/$CURRENT_WORKFLOW_LOG_NAME.txt"
echo -e "$(set_color "$(set_color "  stop " "$BOLD_INTNS_WHITE")" "$INTNS_BG_BLUE") Pipeline ran for $((duration / 60)) minutes and $((duration % 60)) seconds."

echo ""
echo -e $(set_color "Logs are available at $CURRENT_WORKFLOW_LOG_NAME" "$BOLD_INTNS_RED")

if [ $NOTIFICATIONS_ENABLED -eq 1 ]; then 
    EMAIL=$(create_workflow_notification_email "$NOTIFICATION_TITLE" "$NOTIFICATION_BODY")
    send_email "$EMAIL_TO" "Bash CI/CD completion" "$EMAIL" "$(get_secret key resend_key)"
fi

exit 0