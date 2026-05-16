-- Migration: 130_drop_channel_monitor_tables
-- 前向删除已下线的渠道监控模块表；保留历史 migration 不回写，避免破坏已部署环境迁移链。

DROP TABLE IF EXISTS channel_monitor_daily_rollups;
DROP TABLE IF EXISTS channel_monitor_histories;
DROP TABLE IF EXISTS channel_monitors;
DROP TABLE IF EXISTS channel_monitor_request_templates;
