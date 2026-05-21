#!/usr/bin/env bash

send_email()
{
    local recipient=$1 subject=$2 html_content=$3

    jq -n \
    --arg from "$FRIENDLY_NAME <$EMAIL_FROM>" \
    --arg to "$recipient" \
    --arg subject "$subject" \
    --arg html "$html_content" '
    {
    "from": $from,
    "to": [$to],
    "subject": $subject,
    "html": $html
    }' | curl -X POST https://api.resend.com/emails \
        -H "Authorization: Bearer $RESEND_KEY" \
        -H "Content-Type: application/json" \
        -d @-
}

# array for template parts for later concatenation
declare -A email_parts=(
    [header_and_title_partial]=
    [results_header]=
    [footer]=
)

# parts of the template are assigned below
IFS='' read -r -d '' email_parts["footer"] <<"EOF"
</table>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding:40px 30px; background: linear-gradient(20deg, #5d2727, #b65252); color:white; text-align:center;">
                            <h3 style="font-size:24px; margin:0 0 15px 0;">Like the project?</h3>
                            <p style="font-size:18px; margin:0 0 25px 0; line-height:1.5;">
                                Throw a star at the repo or hop on board as a contributor
                            </p>
                            <a href="https://github.com/mrtnstl/bash-ci-cd"  target="_blank"
                               style="background-color:white; color:#212222; padding:16px 36px; text-decoration:none; font-weight:bold; font-size:17px; border-radius:6px; display:inline-block;">
                                CHECK OUT THE PROJECT
                            </a>

                        </td>
                    </tr>
                    
                    <tr>
                        <td style="background-color:#202124; color:#9aa0a6; padding:16px 30px; text-align:center; font-size:13px;">
                            <p style="margin:0 0 0px 0;">
                                This email was produced by <strong><em>Bash CI/CD</em></strong>
                            </p>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
    
</body>
</html>
EOF

read -r -d '' email_parts["results_header"] <<"EOF"
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="padding:35px 30px; background-color:#222222; text-align: left; border-radius: 12px;">
    <tr>
        <td>
            <h4 style="margin:0 0 10px 0; color:#cdcdcd;">Tasks</h4>
        </td>
    </tr>
EOF

read -r -d '' email_parts["header_and_title_partial"] <<"EOF"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" type="text/css" href="http://fonts.googleapis.com/css?family=Ubuntu:regular,bold&subset=Latin">

    <title>Bash CI/CD workflow notification</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
    
    <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="background-color:#f4f4f4;">
        <tr>
            <td align="center">
                
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" style="max-width:600px; background-color:#2b2b2b; border-collapse:collapse;">
                     <tr>
                        <td style="font-family: 'Ubuntu'; text-align: center;  padding:30px 10px 8px 10px; color:#fa4c4c; line-height: 21px; letter-spacing: 1px; word-spacing: 5px">
                            <p style="margin:0px 0 0 0; font-size:18px; white-space: preserve nowrap;"> ______     ______     ______     __  __   </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">/\  __ \   /\  __ \   /\  ___\   /\ \_\ \  </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">\ \  __<   \ \  __ \  \ \___  \  \ \  __ \ </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;"> \ \_____\  \ \_\ \_\  \/\_____\  \ \_\ \_\</p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">  \/_____/   \/_/\/_/   \/_____/   \/_/\/_/</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Ubuntu'; text-align: center;  padding:0px 10px 30px 10px; color:#fa4c4c; line-height: 21px; letter-spacing: -1.5px; word-spacing: 5px">
                            <p style="margin:0px 0 0 0; font-size:18px; white-space: preserve nowrap;"> ______     __        ______     _____    </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">/\  ___\   /\ \      /\  ___\   /\  __ \  </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">\ \ \____  \ \ \     \ \ \____  \ \ \_\ \ </p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;"> \ \_____\  \ \_\     \ \_____\  \ \_____|</p>
                            <p style="margin:0px 0 0 0; font-size:20px; white-space: preserve nowrap;">  \/_____/   \/_/      \/_____/   \/____/ </p>
                        </td>
                    </tr>
                    <tr>
                        <td style="background-color:#b65252; padding:30px 20px; text-align:center; color:white;">
                            <h1 style="margin:0; font-size:28px; font-weight:bold;">Workflow Report &#128202;</h1>
                            <p style="margin:10px 0 0 0; font-size:18px; opacity: 0.7">by @mrtnstl</p>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding:40px 30px; background-color:#303030; text-align:left;">
EOF

unset IFS