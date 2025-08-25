<div align="center">

# سیستم مدیریت مرخصی (Leave Tracker)

[![Go Version](https://img.shields.io/github/go-mod/go-version/goravel/framework)](https://go.dev/)
[![License](https://img.shields.io/github/license/goravel/framework)](LICENSE)

</div>

<div align="right" dir="rtl">

## معرفی

سیستم مدیریت مرخصی (Leave Tracker) یک سرویس RESTful API برای مدیریت درخواست‌های مرخصی کارکنان است. این سیستم بر اساس فریم‌ورک Goravel (نسخه Go) توسعه یافته و امکان ثبت، مشاهده و حذف درخواست‌های مرخصی را فراهم می‌کند.

### ویژگی‌های کلیدی

- ثبت درخواست مرخصی (روزانه و ساعتی)
- مشاهده تاریخچه مرخصی‌های کاربر
- گزارش‌گیری از مرخصی‌ها
- احراز هویت کاربران
- اعتبارسنجی پیشرفته درخواست‌ها
- مدیریت رویدادها برای ارسال اعلان

## پیش‌نیازها

- Go 1.16 یا بالاتر
- پایگاه داده (MySQL/PostgreSQL/SQLite)
- (اختیاری) Redis برای مدیریت صف‌ها

## نصب و راه‌اندازی

### 1. کپی فایل تنظیمات

```bash
cp .env.example .env
```

### 2. تنظیم متغیرهای محیطی

فایل `.env` را با مقادیر مناسب ویرایش کنید:

```env
APP_NAME=LeaveTracker
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_URL=http://localhost:8080

# تنظیمات پایگاه داده
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=leave_tracker
DB_USERNAME=root
DB_PASSWORD=

# تنظیمات احراز هویت JWT
JWT_SECRET=your-jwt-secret-key
```

### 3. نصب وابستگی‌ها

```bash
go mod download
```

### 4. تولید کلید برنامه

```bash
go run . artisan key:generate
```

### 5. اجرای مایگریشن‌ها

```bash
go run . artisan migrate
```

### 6. اجرای سرور توسعه

```bash
go run . artisan serve
```



</div>
