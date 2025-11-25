-- 部署マスタ初期データ
INSERT INTO departments (code, name, parent_id, manager_employee_number) VALUES
('DEV', '開発部', NULL, NULL),
('SALES', '営業部', NULL, NULL),
('HR', '人事部', NULL, NULL);

-- 役職マスタ初期データ
INSERT INTO positions (code, name, level) VALUES
('MGR', 'マネージャー', 5),
('LEAD', 'リーダー', 4),
('SR', 'シニア', 3),
('MID', 'ミドル', 2),
('JR', 'ジュニア', 1);

-- テスト用従業員データ
INSERT INTO employees (name, name_kana, birthday, gender_cd, email, department_id, position_id, hire_date, employment_type) VALUES
('山田太郎', 'ヤマダタロウ', '1990-01-01', 1, 'yamada@example.com', 1, 3, '2020-04-01', 1),
('佐藤花子', 'サトウハナコ', '1992-05-15', 2, 'sato@example.com', 2, 2, '2021-04-01', 1);

-- テスト用ユーザーアカウント（パスワード: password123）
-- bcrypt hash of "password123"
INSERT INTO users (employee_number, email, password_hash, role) VALUES
(1, 'yamada@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
(2, 'sato@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'employee');
