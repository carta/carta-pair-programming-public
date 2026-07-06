import { Routes, Route } from "react-router-dom";
import UsersPage from "@/pages/UsersPage";
import AccountsPage from "@/pages/AccountsPage";
import TransactionsPage from "@/pages/TransactionsPage";

function App() {
  return (
    <Routes>
      <Route path="/" element={<UsersPage />} />
      <Route path="/users/:userId/accounts" element={<AccountsPage />} />
      <Route
        path="/users/:userId/accounts/:accountId/transactions"
        element={<TransactionsPage />}
      />
    </Routes>
  );
}

export default App;
