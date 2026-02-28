import { useList, useCreate, useDelete, useCustom } from "@refinedev/core";
import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

interface ManagedCollection {
  id: string;
  collection_name: string;
}

interface PBCollection {
  id: string;
  name: string;
  type: string;
}

export function ManagedCollectionList() {
  const [confirmAction, setConfirmAction] = useState<{
    type: "register" | "unregister";
    name: string;
    id?: string;
  } | null>(null);

  const { query: managedQuery, result: managedResult } = useList<ManagedCollection>({
    resource: "iam_managed_collections",
  });
  const managed = managedResult.data ?? [];

  const { query: collectionsQuery } = useCustom<PBCollection[]>({
    url: "/api/collections",
    method: "get",
  });
  const allCollectionsRaw = collectionsQuery.data?.data;
  const allCollections: PBCollection[] = Array.isArray(allCollectionsRaw)
    ? allCollectionsRaw
    : (allCollectionsRaw as unknown as { items?: PBCollection[] })?.items ?? [];

  const managedNames = new Set(managed.map((m: ManagedCollection) => m.collection_name));
  const available = allCollections.filter(
    (c) => !managedNames.has(c.name) && !c.name.startsWith("iam_") && c.name !== "_superusers"
  );

  const { mutate: createRecord } = useCreate();
  const { mutate: deleteRecord } = useDelete();

  function handleRegister(name: string) {
    createRecord(
      { resource: "iam_managed_collections", values: { collection_name: name } },
      { onSuccess: () => setConfirmAction(null) }
    );
  }

  function handleUnregister(id: string) {
    deleteRecord(
      { resource: "iam_managed_collections", id },
      { onSuccess: () => setConfirmAction(null) }
    );
  }

  if (managedQuery.isLoading) return <p>Loading...</p>;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="mb-4 text-2xl font-bold">Managed Collections</h1>
        <p className="mb-4 text-sm text-muted-foreground">
          IAM enforcement is applied to managed collections. Register a collection to enable
          permission checks on its records.
        </p>
      </div>

      {/* Currently Managed */}
      <div>
        <h2 className="mb-3 text-lg font-semibold">Currently Managed</h2>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Collection Name</TableHead>
                <TableHead className="w-32">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {managed.length ? (
                managed.map((m: ManagedCollection) => (
                  <TableRow key={m.id}>
                    <TableCell className="font-mono text-sm">{m.collection_name}</TableCell>
                    <TableCell>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() =>
                          setConfirmAction({ type: "unregister", name: m.collection_name, id: m.id })
                        }
                      >
                        Unregister
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={2} className="h-16 text-center">
                    No managed collections.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Available Collections */}
      <div>
        <h2 className="mb-3 text-lg font-semibold">Available Collections</h2>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Collection Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead className="w-32">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {available.length ? (
                available.map((c) => (
                  <TableRow key={c.id}>
                    <TableCell className="font-mono text-sm">{c.name}</TableCell>
                    <TableCell className="text-sm text-muted-foreground">{c.type}</TableCell>
                    <TableCell>
                      <Button
                        size="sm"
                        onClick={() => setConfirmAction({ type: "register", name: c.name })}
                      >
                        Register
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={3} className="h-16 text-center">
                    No available collections to register.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Confirmation Dialog */}
      <AlertDialog open={!!confirmAction} onOpenChange={() => setConfirmAction(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {confirmAction?.type === "register" ? "Register Collection" : "Unregister Collection"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {confirmAction?.type === "register"
                ? `This will enable IAM enforcement on "${confirmAction.name}". All authenticated requests to this collection will be checked against IAM policies.`
                : `This will disable IAM enforcement on "${confirmAction?.name}". Requests will no longer be checked against IAM policies.`}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (confirmAction?.type === "register") {
                  handleRegister(confirmAction.name);
                } else if (confirmAction?.id) {
                  handleUnregister(confirmAction.id);
                }
              }}
            >
              {confirmAction?.type === "register" ? "Register" : "Unregister"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
