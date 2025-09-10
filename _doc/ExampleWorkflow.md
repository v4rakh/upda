# Example Workflow

1. **Deployment of upda**
    - You install and run _upda_ on a central machine (e.g., a monitoring server).

2. **Configure a Source of Updates**
    - Suppose you are managing several Docker containers.
    - You set up **diun** (Docker Image Update Notifier) on one of your hosts.
    - diun is configured to send update notifications via *webhook* to upda’s endpoint whenever a new image version is
      detected.

3. **Sending Update Data**
    - A new version of an image (e.g., `nginx:1.27.1`) becomes available.
    - diun detects this and calls the **upda webhook** with data such as:
        - Host: `node-01`
        - Application/Image: `nginx`
        - Current version: `1.27.0`
        - New version: `1.27.1`
        - Timestamp and metadata

    - Alternatively, for a custom update, you could simply send data using the CLI:
      ```bash
      upda webhook send '{"host": "node-02", "application": "myservice", "version": "2.3.4"}'
      ```
4. **upda Processes and Displays the Update**
    - _upda_ receives the incoming update and stores it.
    - In the dashboard UI, you now see:
        - Which host the update came from
        - The application or container affected
        - The available newer version
        - Update status (pending, approved, etc.)

5. **Managing Update States**
    - From the dashboard, you can mark an update as:
        - *Pending*
        - *Approved*
        - *Ignored*
    - Over time, you get a clear overview of which systems are fully up-to-date and which still need maintenance,
      especially using _upda_'s filters in the overview.

This way, **upda becomes your central hub**, aggregating update notifications from automated tools like diun, custom
scripts, or even manual input, turning them into an actionable and trackable list.
