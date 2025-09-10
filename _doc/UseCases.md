# Use Cases

### Managing Updates (states: pending, ignored, approved)

**Use Case:**

- A sysadmin is responsible for keeping web servers patched.
- New security updates for `nginx` or `postgresql` appear.
- They come into _upda_ as *pending*.
- If an update is deemed important (security critical), it is marked *approved* to act on.
- If the update is for a deprecated service or a development-only host, it can be *ignored*.
- This gives admins a **clear backlog and prioritization flow** for updates.

### Managing Webhooks (retrieving updates from external systems)

**Use Case:**

- A Docker image update notifier (e.g., **diun**) detects a new image version (`redis:7.2`).
- diun calls a unique webhook URL in upda.
- _upda_ records that host `node-02` is running an outdated image.
- For system packages, a cron job or bash script can push package updates to _upda_ by hitting its webhook.
- This allows **any system or tool** (CI/CD pipeline, monitoring tool) to feed update data into one central dashboard.

### Managing Actions (triggering integrations via shoutrrr)

**Use Case:**

- When a critical update is created (e.g., kernel vulnerability), _upda_ triggers a **Slack notification** to the ops team
  using shoutrrr.
- When the update state changes to *approved*, it can trigger a **Jenkins job** to run an Ansible playbook that rolls
  out the patch.
- This integrates updates into the **automation workflow**, ensuring they’re not just tracked but also trigger
  remediation steps.

### History of Actions

**Use Case:**

- An admin investigates whether an alert about a PostgreSQL update was actually sent to the team.
- The action history shows that the Slack message was delivered at 14:32 yesterday.
- This provides **auditability** and helps diagnose gaps in communication or automation.

### Viewing Events (tracking changes to updates)

**Use Case:**

- An admin team member marks `nginx:1.27.1` as approved but later another colleague changes it to ignored.
- The event history clearly shows who made the change and when.
- This is valuable for **collaboration and audit compliance**, ensuring changes are transparent and traceable.

### Metrics Exporter (Prometheus integration)

**Use Case:**

- DevOps teams export metrics to Prometheus and view them in Grafana dashboards.
- They track:
    - Number of *pending* updates by host
    - Trend of how often new updates appear
    - Average time it takes to move updates from *pending* to *completed*
- These metrics help measure **update compliance** and **system freshness**, which are critical for security posture.  