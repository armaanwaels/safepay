CREATE TABLE IF NOT EXISTS invoices (
  id TEXT PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  customer_id TEXT NOT NULL,
  amount INTEGER NOT NULL,
  currency TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'unpaid',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
  id TEXT PRIMARY KEY,
  invoice_id TEXT NOT NULL,
  method TEXT NOT NULL,
  status TEXT NOT NULL,
  paid_at DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);
