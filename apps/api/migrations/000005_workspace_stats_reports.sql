-- Version: 000005
-- Description: Workspace isolation for stats and reports

-- daily_stats workspace isolation
ALTER TABLE daily_stats ADD COLUMN IF NOT EXISTS workspace_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_daily_stats_workspace_date ON daily_stats(workspace_id, date);

-- report_schedules workspace isolation
ALTER TABLE report_schedules ADD COLUMN IF NOT EXISTS workspace_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_report_schedules_workspace ON report_schedules(workspace_id);
