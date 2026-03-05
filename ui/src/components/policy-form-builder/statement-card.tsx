import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { TagListInput } from "./tag-list-input";
import type { PolicyStatement } from "@/types/policy";

interface StatementCardProps {
  index: number;
  statement: PolicyStatement;
  onChange: (statement: PolicyStatement) => void;
  onRemove: () => void;
  canRemove: boolean;
}

export function StatementCard({
  index,
  statement,
  onChange,
  onRemove,
  canRemove,
}: StatementCardProps) {
  return (
    <Card>
      <CardContent className="space-y-3 pt-4">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Statement {index + 1}</span>
          {canRemove && (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onRemove}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1">
            <Label className="text-xs">SID</Label>
            <Input
              value={statement.sid}
              onChange={(e) =>
                onChange({ ...statement, sid: e.target.value })
              }
              placeholder="e.g. AllowReadPosts"
              className="h-8 text-sm"
            />
          </div>

          <div className="space-y-1">
            <Label className="text-xs">Effect</Label>
            <div className="flex gap-4 pt-1">
              <label className="flex items-center gap-1.5 text-sm">
                <input
                  type="radio"
                  name={`effect-${index}`}
                  checked={statement.effect === "Allow"}
                  onChange={() =>
                    onChange({ ...statement, effect: "Allow" })
                  }
                />
                Allow
              </label>
              <label className="flex items-center gap-1.5 text-sm">
                <input
                  type="radio"
                  name={`effect-${index}`}
                  checked={statement.effect === "Deny"}
                  onChange={() =>
                    onChange({ ...statement, effect: "Deny" })
                  }
                />
                Deny
              </label>
            </div>
          </div>
        </div>

        <TagListInput
          label="Actions"
          values={statement.action}
          onChange={(action) => onChange({ ...statement, action })}
          placeholder="e.g. collections:posts:read"
        />

        <TagListInput
          label="Resources"
          values={statement.resource}
          onChange={(resource) => onChange({ ...statement, resource })}
          placeholder="e.g. *"
        />
      </CardContent>
    </Card>
  );
}
