INSERT INTO employees(
    name,
    name_kana,
    age,
    gender_cd,
    birthday
)VALUES(
    '山田太郎',
    'ヤマダタロウ',
    18,
    1,
    cast('2004-01-01' as date)
);

INSERT INTO attendances(
    employee_id ,
    opening_time,
    closing_time,
    attendance_status
)VALUES(
    '1',
    cast('2030-04-01 10:00:00' as datetime),
    cast('2030-04-01 20:00:00' as datetime),
    1,
);
