// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analytics
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

pipeline {

        agent {
                label 'microagent'
        }

        post {

                failure {
                        echoFailure("Pipeline Failed.")
                        slackFailure(
                                'channel',
                                'Harbor Analyst (Builder)',
                                "Build #${env.BUILD_NUMBER} for building harbor-analyst failed." +
                                " Build: ${env.BUILD_URL}"
                        )
                }

                success {
                        echoSuccess("Pipeline Succeeded.")
                        slackSuccess(
                                'channel',
                                'Harbor Analyst (Builder)',
                                "Build #${env.BUILD_NUMBER} for building harbor-analyst succeeded." +
                                " Build: ${env.BUILD_URL}"
                        )
                }

        }

        stages {

                stage('Build Data Fetcher and Analyst (Parallel)') {

                        steps {

                                parallel(
                                        build_data_fetcher: {
                                                sh 'make build-raw-data-fetcher'
                                        },
                                        build_analyst: {
                                                sh 'make build'
                                        }
                                )

                        }

                }

        }

}
