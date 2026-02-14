import { ListTransactionsResponse } from "@/types/transaction";
import { apiFetch } from "./client";

export async function listTransactions(): Promise<ListTransactionsResponse> {
  return apiFetch<ListTransactionsResponse>("/api/v1/transactions");
}
