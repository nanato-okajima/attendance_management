CREATE TABLE attendance_management.employees(
    employee_number int UNSIGNED AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    name_kana varchar(255) NOT NULL,
    age TINYINT UNSIGNED NOT NULL,
    gender_cd TINYINT UNSIGNED NOT NULL,
    birthday DATE NOT NULL,
    PRIMARY KEY (employee_number)
);

CREATE TABLE attendance_management.attendances(
    attendance_id int UNSIGNED AUTO_INCREMENT,
    employee_id int UNSIGNED,
    opening_time DATETIME NOT NULL,
    closing_time DATETIME NOT NULL,
    attendance_status TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (attendance_id),
    FOREIGN KEY (employee_id) REFERENCES employees(employee_number) ON DELETE CASCADE
);
