# alertmanager-msteams

`alertmanager-msteams` is a lightweight Go-based service that forwards Prometheus Alertmanager notifications to Microsoft Teams channels using Adaptive Cards. This tool enables seamless integration between Prometheus alerts and Microsoft Teams, ensuring that your team stays informed about critical events in your infrastructure.

## Features

- **Adaptive Cards Support**: Uses Microsoft's Adaptive Cards for rich and interactive message formatting in Teams channels.
- **Lightweight and Efficient**: Built with Go, providing a statically compiled binary for easy deployment without external dependencies.
- **Simple Configuration**: Minimal setup required—just point Alertmanager to this service.

## Prerequisites

- **Microsoft Teams Webhook URL**: Set up an incoming webhook in your desired Teams channel to obtain the webhook URL.
- **Go Environment**: Ensure you have Go installed if you plan to build from source.
- **Docker**: For containerized deployment.

## Installation

### Using Docker

1. **Pull the Docker image**:

   ```sh
   docker pull ghcr.io/karunsiri/alertmanager-msteams:latest
   ```

2. **Run the container**:

   ```sh
   docker run -d -p 8080:8080 \
       -e WEBHOOK_URL="https://outlook.office.com/webhook/your-webhook-url" \
       --name alertmanager-msteams \
       ghcr.io/karunsiri/alertmanager-msteams:latest
   ```

   Replace `https://outlook.office.com/webhook/your-webhook-url` with your actual Teams webhook URL.

### Building from Source

1. **Clone the repository**:

   ```sh
   git clone https://github.com/karunsiri/alertmanager-msteams.git
   cd alertmanager-msteams
   ```

2. **Build the binary**:

   ```sh
   go build -o alertmanager-msteams
   ```

3. **Run the application**:

   ```sh
   WEBHOOK_URL="https://outlook.office.com/webhook/your-webhook-url" ./alertmanager-msteams
   ```

## Configuration

The service can be configured using environment variables:

- `WEBHOOK_URL` (Required) – The Microsoft Teams incoming webhook URL.
- `PORT` (Optional) – Port on which the service listens. Defaults to `8080`.

## Usage

1. **Configure Prometheus Alertmanager**:

   In your Alertmanager configuration file (`alertmanager.yml`), add a webhook receiver pointing to the `alertmanager-msteams` service:

   ```yaml
   receivers:
     - name: 'msteams'
       webhook_configs:
         - url: 'http://localhost:8080/alert'
   ```

   Adjust the URL according to your deployment setup.

## Limitations

- **Templates are currently not customizable**. The service uses a built-in format for Teams messages and does not yet support user-defined templates. Future versions may add template customization.
- **Only a single webhook URL is supported**. If you need multiple Teams channels, you will need to run multiple instances of `alertmanager-msteams`.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to enhance the functionality or fix bugs.

## Testing

TBA

---

