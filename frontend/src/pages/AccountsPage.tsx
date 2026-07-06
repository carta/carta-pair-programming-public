import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface Account {
  id: number;
  name: string;
  account_type: string;
  notes: string | null;
  balance: string;
}

export default function AccountsPage() {
  const { userId } = useParams<{ userId: string }>();
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`/api/users/${userId}/accounts`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(setAccounts)
      .catch((err) => setError(err.message));
  }, [userId]);

  if (error) return <p className="p-8 text-red-500">Error: {error}</p>;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <Link to="/" className="text-sm text-blue-600 hover:underline">
        &larr; Back to users
      </Link>
      <h1 className="text-2xl font-bold mt-4 mb-6">Accounts</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {accounts.map((account) => (
          <Link
            key={account.id}
            to={`/users/${userId}/accounts/${account.id}/transactions`}
          >
            <Card className="hover:shadow-md transition-shadow cursor-pointer">
              <CardHeader>
                <CardTitle>{account.name}</CardTitle>
                <CardDescription className="capitalize">
                  {account.account_type}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-semibold">
                  ${parseFloat(account.balance).toLocaleString("en-US", {
                    minimumFractionDigits: 2,
                  })}
                </p>
                <p className="text-sm text-gray-500 mt-1">
                  {account.notes?.substring(0, 50)}
                </p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
