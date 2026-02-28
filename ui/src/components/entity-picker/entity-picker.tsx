import { useState } from "react";
import { useList } from "@refinedev/core";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";

interface EntityPickerProps {
  /** PocketBase collection name (e.g., "iam_policies", "users") */
  resource: string;
  /** Field to display as label (e.g., "name", "email") */
  labelField: string;
  /** Field to use as secondary text (optional) */
  secondaryField?: string;
  /** Trigger button text */
  triggerLabel: string;
  /** IDs to exclude from the list (already attached) */
  excludeIds?: string[];
  /** Called when an entity is selected */
  onSelect: (id: string, record: Record<string, unknown>) => void;
}

export function EntityPicker({
  resource,
  labelField,
  secondaryField,
  triggerLabel,
  excludeIds = [],
  onSelect,
}: EntityPickerProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");

  const { query, result } = useList({
    resource,
    filters: search
      ? [{ field: labelField, operator: "contains", value: search }]
      : [],
    pagination: { pageSize: 20 },
    queryOptions: { enabled: open },
  });

  const records = (result.data ?? []).filter(
    (r: Record<string, unknown>) => !excludeIds.includes(r.id as string)
  );

  const handleSelect = (record: Record<string, unknown>) => {
    onSelect(record.id as string, record);
    setOpen(false);
    setSearch("");
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Plus className="mr-1 h-3 w-3" />
          {triggerLabel}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Select {triggerLabel.replace(/^\+ ?/, "")}</DialogTitle>
        </DialogHeader>
        <Input
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          autoFocus
        />
        <div className="max-h-64 overflow-y-auto">
          {query.isLoading ? (
            <p className="py-4 text-center text-sm text-muted-foreground">Loading...</p>
          ) : records.length === 0 ? (
            <p className="py-4 text-center text-sm text-muted-foreground">No results</p>
          ) : (
            <ul className="space-y-1">
              {records.map((record: Record<string, unknown>) => (
                <li key={record.id as string}>
                  <button
                    type="button"
                    className="w-full rounded-md px-3 py-2 text-left text-sm hover:bg-accent"
                    onClick={() => handleSelect(record)}
                  >
                    <span className="font-medium">
                      {record[labelField] as string}
                    </span>
                    {secondaryField && !!record[secondaryField] && (
                      <span className="ml-2 text-muted-foreground">
                        {record[secondaryField] as string}
                      </span>
                    )}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
