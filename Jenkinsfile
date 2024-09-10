pipeline {
  agent {
    label 'jenkins-agent-be'
  }

  tools {
    go 'go-1.23.1'
  }

  environment {
    REGION = 'asia-southeast2'
    IMAGE_NAME = 'app'
    REPOSITORY = 'invoice-service'
    PROJECT_ID = credentials('gcp-project-id')
    GOOGLE_CREDENTIALS = credentials('google-service-account-key')
    GOLANGCI_LINT_VERSION = 'v1.61.0'
  }

  stages {
    stage('Set Target Branch') {
      steps {
        script {
          echo "GIT_BRANCH: ${env.GIT_BRANCH}"
          if (env.GIT_BRANCH == 'origin/master') {
            env.TARGET_BRANCH = 'master'
          } else if (env.GIT_BRANCH == 'origin/development') {
            env.TARGET_BRANCH = 'development'
          }
        }
      }
    }

    stage('Checkout Code') {
      steps {
        script {
          def repoUrl = 'https://github.com/faisalilhami27/invoice-service.git'

          checkout([$class: 'GitSCM',
            branches: [
              [name: "*/${env.TARGET_BRANCH}"]
            ],
            userRemoteConfigs: [
              [url: repoUrl, credentialsId: 'github-credential']
            ]
          ])
        }
      }
    }

    stage('Install Dependencies') {
      steps {
        script {
          sh '''
            export GOPATH=$(go env GOPATH)
            export PATH=$PATH:$GOPATH/bin

            if ! [ -x "$(command -v golangci-lint)" ]; then
              echo "golangci-lint not found, installing..."
              go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
            else
              echo "golangci-lint is already installed"
            fi
          '''
          sh 'go mod tidy'
        }
      }
    }

    stage('Run Linter') {
      steps {
        sh '''
          export GOPATH=$(go env GOPATH)
          export PATH=$PATH:$GOPATH/bin
          golangci-lint run --out-format html > golangci-lint.html
        '''
      }
    }

    stage('Run Unit Test') {
      steps {
        script {
          sh 'go test'
        }
      }
    }

    stage('Google Cloud Auth') {
      steps {
        script {
          withCredentials([file(credentialsId: 'google-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
            sh '''
            gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
            gcloud config set project ${PROJECT_ID}
            gcloud auth configure-docker ${REGION}-docker.pkg.dev
            '''
          }
        }
      }
    }

    stage('Build Docker Image') {
      steps {
        script {
          sh 'docker build --platform linux/amd64 -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMAGE_NAME}:${BUILD_NUMBER} .'
        }
      }
    }

    stage('Push Docker Image') {
      steps {
        script {
          sh 'docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMAGE_NAME}:${BUILD_NUMBER}'
        }
      }
    }
  }
}
