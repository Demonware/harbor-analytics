// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

pipeline {

        agent {
                label 'microagent'
        }

        environment {
                SLACK_TOKEN = credentials('harbor-analytics-slack-bot-token')
                AWS_ACCESS_KEY_ID = credentials('aws_id_cred_id')
                AWS_SECRET_ACCESS_KEY = credentials('aws_secret_cred_id')
        }

        triggers {
                //Execute the runner pipeline every Friday
                //at around 3pm utc.
                cron('H 15 * * 5')
        }

        post {

                failure {
                        echoFailure("Pipeline Failed.")
                        slackFailure(
                                'channel',
                                'Harbor Analyst (Runner)',
                                "Build #${env.BUILD_NUMBER} for running harbor-analyst failed." +
                                " Build: ${env.BUILD_URL}"
                        )
                }

        }

        stages {

                stage('Pull latest harbor CSVs') {
                        steps {
                                sh 'make get-raw-data'
                        }
                }

                stage('Generate analytics report') {

                        steps {
                                sh 'make run'
                        }

                }

                stage('Publish Report') {
                        steps {
                                sh 'make publish'
                        }
                }

        }

}
