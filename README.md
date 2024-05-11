## Slack to Google Chat Translator

This application translates messages from Slack webhook requests to Google Chat format and forwards them to a Google Chat webhook.

### Build
1. Clone this repository.
2. Navigate to the project directory.
3. Build the Go application using the following command:
    ```bash
    go build -o app
    ```

### Usage
1. Deploy this application and expose it to the internet.
2. Set up a Google Chat webhook.
3. Use the Google Chat webhook, but replace the hostname with your server's hostname. For example:
   - Original: `https://chat.googleapis.com/AAAA/bbb`
   - Replace `https://chat.googleapis.com` with `https://slack-to-google-chat.example.com`
   - Result: `https://slack-to-google-chat.example.com/AAAA/bbb`

### Docker
1. Build the Docker image:
    ```bash
    docker build -t slack-to-google-chat .
    ```

2. Run the Docker container:
    ```bash
    docker run -p 8080:8080 slack-to-google-chat
    ```

### Configuration
- Set the environment variable `PORT` to specify the port on which the server should listen. Default is `8080`.
- Set the list of allowed google chat space ids in the environment variable `ALLOWED_SPACE_IDS`.
