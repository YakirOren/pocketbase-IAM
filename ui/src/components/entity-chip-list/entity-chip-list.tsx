import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";
import { EntityPicker } from "@/components/entity-picker";

interface ChipItem {
  id: string;
  label: string;
}

interface EntityChipListProps {
  /** Section title (e.g., "Attached Policies") */
  title: string;
  /** Currently attached items */
  items: ChipItem[];
  /** PocketBase collection to pick from */
  resource: string;
  /** Field to display as label in picker */
  labelField: string;
  /** Secondary display field in picker */
  secondaryField?: string;
  /** Trigger button text (e.g., "Attach Policy") */
  addLabel: string;
  /** Called when an entity is added */
  onAdd: (id: string, record: Record<string, unknown>) => void;
  /** Called when a chip is removed */
  onRemove: (id: string) => void;
  /** Loading state */
  isLoading?: boolean;
}

export function EntityChipList({
  title,
  items,
  resource,
  labelField,
  secondaryField,
  addLabel,
  onAdd,
  onRemove,
  isLoading = false,
}: EntityChipListProps) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-medium text-muted-foreground">{title}</h3>
      <div className="flex flex-wrap items-center gap-2">
        {isLoading ? (
          <span className="text-sm text-muted-foreground">Loading...</span>
        ) : items.length === 0 ? (
          <span className="text-sm text-muted-foreground">None</span>
        ) : (
          items.map((item) => (
            <Badge key={item.id} variant="secondary" className="gap-1 pr-1">
              {item.label}
              <button
                type="button"
                onClick={() => onRemove(item.id)}
                className="ml-1 rounded-full p-0.5 hover:bg-muted"
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))
        )}
        <EntityPicker
          resource={resource}
          labelField={labelField}
          secondaryField={secondaryField}
          triggerLabel={addLabel}
          excludeIds={items.map((i) => i.id)}
          onSelect={onAdd}
        />
      </div>
    </div>
  );
}
