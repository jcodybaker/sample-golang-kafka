## Getting Started

We provide a sample app that produces and consumes messages to/from Kafka using Golang that you can deploy on App Platform. These steps will get this sample application running for you using App Platform.

**Note: Following these steps may result in charges for the use of DigitalOcean services.**

### Requirements

* You need a DigitalOcean account. If you don't already have one, you can sign up at https://cloud.digitalocean.com/registrations/new.
* You need a running Kafka instance. If you don't already have one, you can create one at https://cloud.digitalocean.com/databases.

## Deploying the App

Click this button to deploy the app to the DigitalOcean App Platform. If you are not logged in, you will be prompted to log in with your DigitalOcean account.

[![Deploy to DigitalOcean](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/blesswinsamuel/sample-golang-kafka/tree/main)

Using this button disables the ability to automatically re-deploy your app when pushing to a branch or tag in your repository as you are using this repo directly.

After clicking the **Deploy to DigitalOcean** button or completing the instructions above to fork the repo, follow these steps:

1. Configure the environment variables to point to your Kafka instance.
2. Provide a name for your app and select which region you want to deploy your app to and click **Next**. The closest region to you should be selected by default. All App Platform apps are routed through a global CDN. So this will not affect your app performance, unless it needs to talk to external services.
3. On the following screen, leave all the fields as they are and click **Next**.
4. Confirm your **Plan** settings and how many containers you want to launch and click **Launch Basic/Pro App**.
5. You should see a "Building..." progress indicator. You can click **View Logs** to see more details of the build.
6. It can take a few minutes for the build to finish, but you can follow the progress in the **Deployments** tab.
7. Once the build completes successfully, right click on the **Live App** link in the header, click "Copy Link Address" and run `curl -X POST https://xxx.ondigitalocean.app/produce --data "hello world!"` in a terminal. Go to the Runtime Logs of the consumer and you should now see message consumed logs.

If you want to automatically re-deploy your app, [fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo) the GitHub repository to your account so that you have a copy of it stored to the cloud. Click the **Fork** button in the GitHub repository and follow the on-screen instructions.

After forking the repo, you should now be viewing this README in your own GitHub org (e.g. `https://github.com/<your-org>/sample-golang-kafka`). To deploy the new repo, make a couple of changes to the `.do/app.yaml` file.

1. Update the environment variables under `envs` -> `value` to point to your Kafka instance.
2. Set your repo path under `services` -> `github` -> `repo`.

Once the above changes are made, run `doctl apps create --spec .do/app.yaml`.

1. Go to https://cloud.digitalocean.com/apps, and select your app. You should see a "Building..." progress indicator. You can click **View Logs** to see more details of the build.
1. It can take a few minutes for the build to finish, but you can follow the progress in the **Deployments** tab.
1. Once the build completes successfully, right click on the **Live App** link in the header, click "Copy Link Address" and run `curl -X POST https://xxx.ondigitalocean.app/produce --data "hello world!"` in a terminal. Go to the Runtime Logs of the consumer and you should now see message consumed logs.

### Making Changes to Your App

If you followed the steps to fork the repo and used your own copy when deploying the app, you can push changes to your fork and see App Platform automatically re-deploy the update to your app. During these automatic deployments, your application will never pause or stop serving request because App Platform offers zero-downtime deployments.

### Learn More

You can learn more about the App Platform and how to manage and update your application at https://www.digitalocean.com/docs/app-platform/.

## Deleting the App

When you no longer need this sample application running live, you can delete it by following these steps:
1. Visit the Apps control panel at https://cloud.digitalocean.com/apps.
2. Navigate to the sample app.
3. In the **Settings** tab, click **Destroy**.

**Note: If you do not delete your app, charges for using DigitalOcean services will continue to accrue.**
