# AWS Bedrock Spend & Claude Usage by Account

**Window:** past 3 months
**Source dashboard:** [AWS Bedrock Spend](https://app.datadoghq.com/dashboard/w8w-zrn-qwc)
**Org cross-reference:** `terraform/aws/carta-organization/main.tf` (`aws_organizations_account` resources)

Every AWS account visible in Datadog Cloud Cost over the window is listed, including those with **$0 Bedrock spend** and **0 Claude invocations**.

- **Account universe** = `aws.cost.net.amortized.shared.resources.allocated{*}` grouped by `aws_member_account_name`,`aws_member_account_id` (any account with any AWS cost in the window).
- **Bedrock spend** = same metric filtered to `aws_product:*bedrock*`, Usage/DiscountedUsage/SavingsPlanCoveredUsage, excl. enterprise support. `$0.00` = no Bedrock cost.
- **Claude invocations** = `carta.bedrock.invocations{model_id:*claude* OR model_id:*anthropic*}` grouped by `account`. App instrumentation only exists for Carta-instrumented accounts — `0` means no instrumentation data, not necessarily zero traffic (e.g. Accelex/Avantia emit Bedrock cost but not this metric).
- **In org TF** = whether the account is managed in `carta-organization/main.tf`. `Yes` = active resource; `Removed` = decommissioned via `removed{}` block (SRE-9834, closed outside TF); `No` = not managed by this stack.

| Account name | Account ID | Bedrock spend (3mo) | Claude invocations (3mo) | In org TF |
|---|---|---:|---:|---|
| acx-amers | 733697723310 | $9,402.86 | 0 | Yes |
| acx-aztec | 816069169295 | $5,381.31 | 0 | Yes |
| acx-clearwater | 911931247427 | $297.34 | 0 | Yes |
| acx-data-science | 014498663992 | $60.00 | 0 | Yes |
| acx-dev | 536930143272 | $21,380.66 | 0 | Yes |
| acx-emea | 888294021526 | $25,756.29 | 0 | Yes |
| acx-factset | 339713127011 | $250.74 | 0 | Yes |
| acx-s3-backups | 757886745538 | $0.00 | 0 | Yes |
| Audit | 367345877668 | $0.00 | 0 | Yes |
| Avantia-Dev | 533267317066 | $118.61 | 0 | Yes |
| Avantia-Kindergarten | 289109111094 | $0.00 | 0 | Yes |
| Avantia-Prod | 774305582328 | $479.28 | 0 | Yes |
| Avantia-Sandbox | 024848487132 | $0.00 | 0 | Yes |
| Avantia-Staging | 149536482438 | $1.91 | 0 | Yes |
| carta-aztec-eu | 588301174527 | $0.00 | 0 | Yes |
| carta-bedrock | 559050237467 | $138,286.21 | 854,914 | Yes |
| carta-clearwater-production | 698387659152 | $0.32 | 0 | Yes |
| carta-data-science | 959229825965 | $30,170.27 | 100,639 | Yes |
| carta-dev | 360845155482 | $8,255.00 | 83,958 | Yes |
| carta-factset-production | 275144393957 | $0.00 | 0 | Yes |
| carta-glean | 011528299920 | $1.24 | 0 | Yes |
| carta-management | 723427871053 | $185.31 | 4,658 | Yes |
| carta-n8n | 383976104849 | $17.94 | 0 | Yes |
| carta-organization | 754112950986 | $0.00 | 0 | Yes (root) |
| carta-people-operations | 686326297373 | $0.00 | 0 | Yes |
| carta-production | 068134724467 | $307,992.47 | 2,353,675 | Yes |
| carta-production-eu | 738497419789 | $1,214.50 | 22,242 | Yes |
| carta-sandbox | 736808520331 | $69,821.63 | 707,686 | Yes |
| carta-snowplow | 263325493398 | $0.00 | 0 | Yes |
| carta-state-street-production-eu | 303613132540 | $42.57 | 187 | Yes |
| carta-test | 936501470813 | $0.00 | 0 | Yes |
| carta-wiz | 296601444025 | $0.00 | 0 | Yes |
| carta-workspaces | 308739448332 | $0.00 | 0 | Yes |
| carta-yearend-dev | 387346179764 | $0.00 | 0 | Yes |
| carta-yearend-prod | 434585959722 | $0.00 | 0 | Yes |
| hris-production | 360390278422 | $0.00 | 0 | Yes |
| hris-test | 258271782950 | $0.00 | 0 | Yes |
| listalpha-dev | 681136933993 | $0.00 | 0 | No |
| listalpha-infra | 539247495207 | $0.00 | 0 | Removed |
| listalpha-mgmt | 224381535716 | $0.00 | 0 | Removed |
| listalpha-stage | 008447347028 | $0.00 | 0 | No |
| Log archive | 358009859755 | $0.00 | 0 | Yes |
| vauban-demo | 579883826072 | $0.00 | 0 | No |
| vauban-dev | 125579976307 | $0.00 | 0 | No |
| vauban-preproduction | 002440740839 | $0.00 | 0 | No |
| vauban-production | 748273700665 | $0.00 | 0 | Yes |
| vauban-sandbox | 100954143592 | $0.00 | 0 | No |

**Totals:** 47 accounts with AWS cost · 21 with Bedrock spend (~$619K) · 8 with Claude instrumentation (~4.13M invocations).

## Cross-reference notes

`carta-organization/main.tf` defines **39 active** `aws_organizations_account` resources + **2 removed** (decommissioned). All 39 active accounts appear in Datadog cost — no managed account is missing.

**8 Datadog accounts are not actively managed by this TF stack** (all have $0 Bedrock spend):

- **`listalpha-infra`, `listalpha-mgmt`** — in TF as `removed{}` blocks (SRE-9834, closed via the AWS Organizations runbook, not TF). They still show AWS cost in the window — expected residual/closing cost, but worth confirming they are fully closed.
- **`listalpha-dev`, `listalpha-stage`** — not present in `carta-organization` TF at all, yet still incurring cost. Possible leftovers from the ListAlpha decommission (SRE-9834) or managed elsewhere.
- **`vauban-dev`, `vauban-sandbox`, `vauban-demo`, `vauban-preproduction`** — only `vauban-production` is managed in this TF; the four non-prod Vauban accounts are not. Likely a separate org/stack or unmanaged.

No Bedrock or Claude activity exists in any unmanaged account, so from a Bedrock standpoint they are all $0 / 0.
