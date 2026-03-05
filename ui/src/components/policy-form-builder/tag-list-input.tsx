import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, X } from "lucide-react";

interface TagListInputProps {
  label: string;
  values: string[];
  onChange: (values: string[]) => void;
  placeholder?: string;
}

export function TagListInput({
  label,
  values,
  onChange,
  placeholder = "Enter value...",
}: TagListInputProps) {
  const updateValue = (index: number, value: string) => {
    const next = [...values];
    next[index] = value;
    onChange(next);
  };

  const removeValue = (index: number) => {
    onChange(values.filter((_, i) => i !== index));
  };

  const addValue = () => {
    onChange([...values, ""]);
  };

  return (
    <div className="space-y-1">
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <div className="flex flex-wrap items-center gap-1.5">
        {values.map((v, i) => (
          <div key={i} className="flex items-center gap-1">
            <Input
              value={v}
              onChange={(e) => updateValue(i, e.target.value)}
              placeholder={placeholder}
              className="h-7 w-56 text-xs"
            />
            {values.length > 1 && (
              <button
                type="button"
                onClick={() => removeValue(i)}
                className="rounded p-0.5 hover:bg-muted"
              >
                <X className="h-3 w-3" />
              </button>
            )}
          </div>
        ))}
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-7 text-xs"
          onClick={addValue}
        >
          <Plus className="mr-1 h-3 w-3" />
          Add
        </Button>
      </div>
    </div>
  );
}
