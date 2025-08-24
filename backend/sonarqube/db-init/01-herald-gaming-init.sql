-- Herald.lol Gaming Analytics - SonarQube Database Initialization
-- PostgreSQL setup for gaming platform code quality analysis

-- Create gaming-specific extensions
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Gaming analytics schemas
CREATE SCHEMA IF NOT EXISTS herald_sonar;
CREATE SCHEMA IF NOT EXISTS gaming_metrics;
CREATE SCHEMA IF NOT EXISTS code_quality;

-- Set proper ownership
ALTER DATABASE herald_sonar OWNER TO herald_sonar;
ALTER SCHEMA herald_sonar OWNER TO herald_sonar;
ALTER SCHEMA gaming_metrics OWNER TO herald_sonar;
ALTER SCHEMA code_quality OWNER TO herald_sonar;

-- Gaming-specific indexes for performance
-- These will be created by SonarQube, but we can optimize for our use case

-- Create gaming metrics table for custom analytics
CREATE TABLE IF NOT EXISTS gaming_metrics.herald_quality_metrics (
    id SERIAL PRIMARY KEY,
    project_key VARCHAR(400) NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    metric_value DECIMAL(30,20),
    gaming_category VARCHAR(100), -- gaming, performance, riot_api, etc.
    analysis_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX(project_key, gaming_category),
    INDEX(analysis_date)
);

-- Gaming performance tracking
CREATE TABLE IF NOT EXISTS gaming_metrics.herald_performance_tracking (
    id SERIAL PRIMARY KEY,
    project_key VARCHAR(400) NOT NULL,
    analysis_duration_ms INTEGER,
    code_lines_analyzed INTEGER,
    gaming_components_analyzed INTEGER, -- LoL/TFT specific components
    performance_target_met BOOLEAN DEFAULT FALSE, -- <5000ms target
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Custom gaming rules violations
CREATE TABLE IF NOT EXISTS code_quality.herald_gaming_violations (
    id SERIAL PRIMARY KEY,
    project_key VARCHAR(400) NOT NULL,
    rule_key VARCHAR(200) NOT NULL,
    gaming_severity VARCHAR(50), -- gaming_critical, gaming_major, etc.
    riot_api_related BOOLEAN DEFAULT FALSE,
    performance_impact VARCHAR(100),
    component_path TEXT,
    violation_count INTEGER DEFAULT 1,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gaming dashboard materialized views for performance
CREATE MATERIALIZED VIEW IF NOT EXISTS gaming_metrics.herald_dashboard_summary AS
SELECT 
    project_key,
    COUNT(*) as total_analyses,
    AVG(analysis_duration_ms) as avg_analysis_time,
    SUM(code_lines_analyzed) as total_lines_analyzed,
    COUNT(*) FILTER (WHERE performance_target_met = true) as performance_compliant_analyses,
    MAX(created_at) as last_analysis
FROM gaming_metrics.herald_performance_tracking
WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY project_key;

-- Index for gaming dashboard performance
CREATE INDEX IF NOT EXISTS idx_herald_dashboard_project_date 
ON gaming_metrics.herald_performance_tracking(project_key, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_gaming_violations_project_severity
ON code_quality.herald_gaming_violations(project_key, gaming_severity);

-- Gaming-specific PostgreSQL optimizations
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET pg_stat_statements.track = 'all';
ALTER SYSTEM SET pg_stat_statements.max = 10000;

-- Gaming database performance tuning
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '16MB';
ALTER SYSTEM SET maintenance_work_mem = '128MB';
ALTER SYSTEM SET max_connections = '200';

-- Gaming-specific logging for debugging
ALTER SYSTEM SET log_min_duration_statement = '1000'; -- Log slow gaming queries
ALTER SYSTEM SET log_line_prefix = '[%t] %u@%d ';

-- Grant permissions for gaming analytics
GRANT ALL PRIVILEGES ON DATABASE herald_sonar TO herald_sonar;
GRANT ALL PRIVILEGES ON SCHEMA herald_sonar TO herald_sonar;
GRANT ALL PRIVILEGES ON SCHEMA gaming_metrics TO herald_sonar;
GRANT ALL PRIVILEGES ON SCHEMA code_quality TO herald_sonar;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gaming_metrics TO herald_sonar;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA code_quality TO herald_sonar;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA gaming_metrics TO herald_sonar;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA code_quality TO herald_sonar;

-- Create gaming-specific functions
CREATE OR REPLACE FUNCTION gaming_metrics.refresh_herald_dashboard()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW gaming_metrics.herald_dashboard_summary;
END;
$$ LANGUAGE plpgsql;

-- Function to track gaming performance metrics
CREATE OR REPLACE FUNCTION gaming_metrics.track_herald_analysis(
    p_project_key VARCHAR(400),
    p_duration_ms INTEGER,
    p_lines_analyzed INTEGER,
    p_gaming_components INTEGER
)
RETURNS void AS $$
DECLARE
    target_met BOOLEAN := (p_duration_ms < 5000); -- Gaming performance target
BEGIN
    INSERT INTO gaming_metrics.herald_performance_tracking (
        project_key,
        analysis_duration_ms,
        code_lines_analyzed,
        gaming_components_analyzed,
        performance_target_met
    ) VALUES (
        p_project_key,
        p_duration_ms,
        p_lines_analyzed,
        p_gaming_components,
        target_met
    );
END;
$$ LANGUAGE plpgsql;

-- Gaming quality metrics aggregation function
CREATE OR REPLACE FUNCTION gaming_metrics.get_herald_quality_summary(
    p_project_key VARCHAR(400),
    p_days INTEGER DEFAULT 30
)
RETURNS TABLE (
    gaming_category VARCHAR(100),
    avg_metric_value DECIMAL(30,20),
    metric_count BIGINT,
    latest_value DECIMAL(30,20)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        hqm.gaming_category,
        AVG(hqm.metric_value) as avg_metric_value,
        COUNT(*) as metric_count,
        LAST_VALUE(hqm.metric_value) OVER (
            PARTITION BY hqm.gaming_category 
            ORDER BY hqm.analysis_date 
            ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
        ) as latest_value
    FROM gaming_metrics.herald_quality_metrics hqm
    WHERE hqm.project_key = p_project_key
      AND hqm.analysis_date >= CURRENT_DATE - INTERVAL '1 day' * p_days
    GROUP BY hqm.gaming_category;
END;
$$ LANGUAGE plpgsql;

-- Set up automatic stats collection for gaming performance
SELECT cron.schedule('herald-stats-collection', '*/10 * * * *', 'SELECT gaming_metrics.refresh_herald_dashboard();');

-- Initial gaming metrics for Herald.lol project
INSERT INTO gaming_metrics.herald_quality_metrics (
    project_key,
    metric_name,
    metric_value,
    gaming_category
) VALUES 
    ('herald-gaming-analytics', 'gaming_performance_target_ms', 5000, 'performance'),
    ('herald-gaming-analytics', 'gaming_concurrent_users_target', 1000000, 'scalability'),
    ('herald-gaming-analytics', 'gaming_uptime_target_percent', 99.9, 'reliability'),
    ('herald-gaming-analytics', 'riot_api_compliance_score', 100, 'compliance')
ON CONFLICT DO NOTHING;

-- Optimize for gaming workload
ANALYZE gaming_metrics.herald_quality_metrics;
ANALYZE gaming_metrics.herald_performance_tracking;
ANALYZE code_quality.herald_gaming_violations;