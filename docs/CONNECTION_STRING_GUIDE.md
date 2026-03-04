# SUPABASE CONNECTION STRING FORMAT GUIDE

## ❌ WRONG - Password with special characters directly:
SUPABASE_DATABASE_URL=postgresql://postgres:MyP@ss#word:123@db.xxxxx.supabase.co:5432/postgres
# This will FAIL because @ and # and : confuse the parser

## ✅ CORRECT - URL-encoded password:
SUPABASE_DATABASE_URL=postgresql://postgres:MyP%40ss%23word%3A123@db.xxxxx.supabase.co:5432/postgres
# Special chars are encoded: @ → %40, # → %23, : → %3A

## 🔧 How to encode:
1. Go to: https://www.urlencoder.org/
2. Paste your password (e.g., "MyP@ss#word:123")
3. Click "Encode"
4. Copy result (e.g., "MyP%40ss%23word%3A123")
5. Replace [PASSWORD] in connection string

## 📝 Common special characters that MUST be encoded:
@ → %40
# → %23
: → %3A
/ → %2F
? → %3F
& → %26
% → %25
= → %3D
+ → %2B
space → %20

## Example .env file:
PORT=8080
ENV=development

# Session mode connection (port 5432)
SUPABASE_DATABASE_URL=postgresql://postgres:EncodedPasswordHere@db.vxbcjtnznelovdeevxur.supabase.co:5432/postgres
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

JWT_SECRET=your-secret-key-here
ALLOWED_ORIGINS=*
