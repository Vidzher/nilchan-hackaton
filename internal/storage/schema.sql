CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL COLLATE NOCASE UNIQUE,
    username TEXT NOT NULL COLLATE NOCASE UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_progress (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    xp INTEGER NOT NULL DEFAULT 0 CHECK (xp >= 0),
    e_points INTEGER NOT NULL DEFAULT 0 CHECK (e_points >= 0),
    current_streak INTEGER NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
    last_completion_at DATETIME,
    avatar_id TEXT NOT NULL DEFAULT 'default_avatar',
    frame_id TEXT NOT NULL DEFAULT 'default_frame',
    title_id TEXT,
    showcase_item_id TEXT
);

CREATE TABLE IF NOT EXISTS user_cosmetics (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL,
    PRIMARY KEY (user_id, item_id)
);

CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('PROCESSING', 'FAILED', 'NOT_COMPLETED', 'COMPLETED')),
    purchased_overflow_slot INTEGER NOT NULL DEFAULT 0 CHECK (purchased_overflow_slot IN (0, 1)),
    created_at DATETIME NOT NULL,
    completed_at DATETIME,
    xp_earned INTEGER CHECK (xp_earned IS NULL OR xp_earned >= 0),
    e_points_earned INTEGER CHECK (e_points_earned IS NULL OR e_points_earned >= 0),
    UNIQUE (user_id, url)
);

CREATE INDEX IF NOT EXISTS idx_resources_user_created_at ON resources(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_resources_user_status ON resources(user_id, status);

CREATE TABLE IF NOT EXISTS resource_tags (
    resource_id INTEGER NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    PRIMARY KEY (resource_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_resource_tags_tag ON resource_tags(tag);

CREATE TABLE IF NOT EXISTS quizzes (
    id INTEGER PRIMARY KEY,
    resource_id INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY,
    quiz_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    correct_answer_id INTEGER NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS answers (
    id INTEGER PRIMARY KEY,
    question_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_questions_quiz_id
    ON questions(quiz_id);
CREATE INDEX IF NOT EXISTS idx_answers_question_id
    ON answers(question_id);
