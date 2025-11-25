-- 従業員マスタ
CREATE TABLE IF NOT EXISTS employees (
    employee_number INT UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    name_kana VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL,
    gender_cd TINYINT UNSIGNED NOT NULL COMMENT '1:男性, 2:女性, 3:その他',
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    department_id INT UNSIGNED,
    position_id INT UNSIGNED,
    hire_date DATE NOT NULL,
    employment_type TINYINT UNSIGNED NOT NULL COMMENT '1:正社員, 2:契約, 3:パート, 4:アルバイト',
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (employee_number),
    INDEX idx_email (email),
    INDEX idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ユーザーアカウント
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'employee') NOT NULL DEFAULT 'employee',
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 出勤記録
CREATE TABLE IF NOT EXISTS attendances (
    attendance_id INT UNSIGNED AUTO_INCREMENT,
    employee_id INT UNSIGNED NOT NULL,
    target_date DATE NOT NULL,
    opening_time DATETIME,
    closing_time DATETIME,
    attendance_status TINYINT UNSIGNED NOT NULL COMMENT '1:出勤, 2:欠勤, 3:遅刻, 4:早退',
    clock_source TINYINT UNSIGNED DEFAULT 1 COMMENT '1:Web, 2:モバイル, 3:LINE',
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    work_hours DECIMAL(5,2),
    overtime_hours DECIMAL(5,2),
    note TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (attendance_id),
    FOREIGN KEY (employee_id) REFERENCES employees(employee_number) ON DELETE CASCADE,
    UNIQUE KEY uk_employee_date (employee_id, target_date),
    INDEX idx_target_date (target_date),
    INDEX idx_employee_date (employee_id, target_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 休暇申請
CREATE TABLE IF NOT EXISTS leave_requests (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    leave_type TINYINT UNSIGNED NOT NULL COMMENT '1:有給, 2:特別休暇, 3:欠勤, 4:半休',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    half_day_type TINYINT UNSIGNED COMMENT '1:午前半休, 2:午後半休',
    reason TEXT NOT NULL,
    approval_status TINYINT UNSIGNED DEFAULT 1 COMMENT '1:承認待ち, 2:承認済み, 3:却下',
    approver_id INT UNSIGNED,
    approved_at DATETIME,
    approval_comment TEXT,
    reject_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_employee (employee_number),
    INDEX idx_status (approval_status),
    INDEX idx_dates (start_date, end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 打刻修正申請
CREATE TABLE IF NOT EXISTS attendance_corrections (
    id INT UNSIGNED AUTO_INCREMENT,
    attendance_id INT UNSIGNED NOT NULL,
    employee_number INT UNSIGNED NOT NULL,
    correction_type TINYINT UNSIGNED NOT NULL COMMENT '1:出勤修正, 2:退勤修正, 3:両方',
    corrected_opening_time DATETIME,
    corrected_closing_time DATETIME,
    reason TEXT NOT NULL,
    approval_status TINYINT UNSIGNED DEFAULT 1 COMMENT '1:承認待ち, 2:承認済み, 3:却下',
    approver_id INT UNSIGNED,
    approved_at DATETIME,
    reject_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id) ON DELETE CASCADE,
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_employee (employee_number),
    INDEX idx_status (approval_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 部署マスタ
CREATE TABLE IF NOT EXISTS departments (
    id INT UNSIGNED AUTO_INCREMENT,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    parent_id INT UNSIGNED,
    manager_employee_number INT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY (manager_employee_number) REFERENCES employees(employee_number) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 役職マスタ
CREATE TABLE IF NOT EXISTS positions (
    id INT UNSIGNED AUTO_INCREMENT,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    level TINYINT UNSIGNED NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- LINE連携ユーザー
CREATE TABLE IF NOT EXISTS line_users (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    line_user_id VARCHAR(255) NOT NULL UNIQUE,
    linked_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    UNIQUE KEY uk_employee (employee_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- LINE連携コード
CREATE TABLE IF NOT EXISTS line_linking_codes (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    code VARCHAR(8) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_expires (expires_at),
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
