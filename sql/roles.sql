-- Connect to database repostats
\c repostats

-- Create Roles & Grant Privileges for Grafana datasource
CREATE USER repostats_readonly WITH PASSWORD '1R_repoSt(ats987';
GRANT CONNECT ON DATABASE repostats TO repostats_readonly;
GRANT USAGE ON SCHEMA public TO repostats_readonly;
GRANT USAGE ON SCHEMA gitee TO repostats_readonly;
GRANT USAGE ON SCHEMA gitee_state TO repostats_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO repostats_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA gitee TO repostats_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA gitee_state TO repostats_readonly;

