-- CyberLab V2 Migration
-- Adds: competitions, teams, container_pool

CREATE TABLE IF NOT EXISTS competitions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS competition_challenges (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    competition_id BIGINT NOT NULL,
    challenge_id BIGINT NOT NULL,
    INDEX idx_competition (competition_id),
    INDEX idx_challenge (challenge_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS competition_scores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    competition_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    team_id BIGINT DEFAULT 0,
    score INT DEFAULT 0,
    solved_count INT DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_comp_user (competition_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS teams (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    captain_id BIGINT,
    score INT DEFAULT 0,
    member_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_captain (captain_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS team_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    team_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_team (team_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
