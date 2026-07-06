import { useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface Category {
  id: number;
  name: string;
}

interface Merchant {
  id: number;
  name: string;
}

interface TransactionAccount {
  id: number;
  name: string;
  account_type: string;
}

interface Transaction {
  id: number;
  date: string;
  amount: string;
  category: Category | null;
  merchant: Merchant | null;
  type: string;
  account: TransactionAccount;
}

const PAGE_SIZE = 20;

export default function TransactionsPage() {
  const { userId, accountId } = useParams<{
    userId: string;
    accountId: string;
  }>();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [loadedKey, setLoadedKey] = useState<string | null>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(null);

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => setDebouncedSearch(search), 300);
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [search]);

  // Loading is derived from whether the in-flight request's key has resolved,
  // so we never set state synchronously inside the fetch effect.
  const requestKey = `${accountId ?? ""}::${debouncedSearch}`;
  const loading = loadedKey !== requestKey;

  useEffect(() => {
    let ignore = false;
    const params = debouncedSearch
      ? `?search=${encodeURIComponent(debouncedSearch)}`
      : "";
    fetch(`/api/accounts/${accountId}/transactions${params}`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (ignore) return;
        setTransactions(data);
        setPage(1);
        setLoadedKey(requestKey);
      })
      .catch((err) => {
        if (!ignore) setError(err.message);
      });
    return () => {
      ignore = true;
    };
  }, [accountId, debouncedSearch, requestKey]);

  if (error) return <p className="p-8 text-red-500">Error: {error}</p>;

  const totalPages = Math.ceil(transactions.length / PAGE_SIZE);
  const start = (page - 1) * PAGE_SIZE;
  const pageTransactions = transactions.slice(start, start + PAGE_SIZE);

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <Link
        to={`/users/${userId}/accounts`}
        className="text-sm text-blue-600 hover:underline"
      >
        &larr; Back to accounts
      </Link>
      <h1 className="text-2xl font-bold mt-4 mb-2">Transactions</h1>
      <p className="text-sm text-gray-500 mb-4">
        {transactions.length} total transactions
      </p>
      <Input
        type="search"
        placeholder="Search by merchant or category..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="mb-4 max-w-sm"
      />

      <div className="rounded-md border bg-white">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Date</TableHead>
              <TableHead>Account</TableHead>
              <TableHead>Merchant</TableHead>
              <TableHead>Category</TableHead>
              <TableHead className="text-right">Amount</TableHead>
              <TableHead>Type</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading
              ? Array.from({ length: PAGE_SIZE }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell className="text-right"><Skeleton className="h-4 w-16 ml-auto" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                  </TableRow>
                ))
              : pageTransactions.map((t) => (
                  <TableRow key={t.id}>
                    <TableCell>{t.date}</TableCell>
                    <TableCell>{t.account.name}</TableCell>
                    <TableCell>{t.merchant?.name ?? "—"}</TableCell>
                    <TableCell>{t.category?.name ?? "—"}</TableCell>
                    <TableCell className="text-right">
                      ${parseFloat(t.amount).toLocaleString("en-US", {
                        minimumFractionDigits: 2,
                      })}
                    </TableCell>
                    <TableCell className="capitalize">{t.type}</TableCell>
                  </TableRow>
                ))}
          </TableBody>
        </Table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-4 mt-4">
          <Button
            variant="outline"
            size="sm"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            Previous
          </Button>
          <span className="text-sm text-gray-600">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
}
