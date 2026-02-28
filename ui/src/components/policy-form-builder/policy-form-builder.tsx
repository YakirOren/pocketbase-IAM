import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Plus } from "lucide-react";
import { StatementCard } from "./statement-card";
import { JsonEditor } from "@/components/json-editor";
import type { JsonEditorHandle } from "@/components/json-editor";
import type { PolicyDocument, PolicyStatement } from "@/types/policy";
import { DEFAULT_STATEMENT } from "@/types/policy";

interface PolicyFormBuilderProps {
  value: PolicyDocument;
  onChange: (doc: PolicyDocument) => void;
}

export function PolicyFormBuilder({ value, onChange }: PolicyFormBuilderProps) {
  const [mode, setMode] = useState<"form" | "json">("form");
  const editorRef = useRef<JsonEditorHandle>(null);
  const [jsonResetKey, setJsonResetKey] = useState(0);

  const updateStatement = (index: number, stmt: PolicyStatement) => {
    const next = [...value.statement];
    next[index] = stmt;
    onChange({ ...value, statement: next });
  };

  const removeStatement = (index: number) => {
    onChange({
      ...value,
      statement: value.statement.filter((_, i) => i !== index),
    });
  };

  const addStatement = () => {
    onChange({
      ...value,
      statement: [...value.statement, { ...DEFAULT_STATEMENT }],
    });
  };

  const handleModeChange = (newMode: string) => {
    if (newMode === "json") {
      // Sync form state to JSON editor
      setJsonResetKey((k) => k + 1);
    } else if (newMode === "form") {
      // Sync JSON editor state back to form
      const raw = editorRef.current?.getValue();
      if (raw) {
        try {
          const parsed = JSON.parse(raw) as PolicyDocument;
          onChange(parsed);
        } catch {
          // Invalid JSON — stay on JSON tab
          return;
        }
      }
    }
    setMode(newMode as "form" | "json");
  };

  return (
    <div className="space-y-4">
      <Tabs value={mode} onValueChange={handleModeChange}>
        <TabsList>
          <TabsTrigger value="form">Form Builder</TabsTrigger>
          <TabsTrigger value="json">Raw JSON</TabsTrigger>
        </TabsList>

        <TabsContent value="form" className="space-y-4">
          <div className="space-y-1">
            <Label className="text-xs">Version</Label>
            <Input
              value={value.version}
              onChange={(e) => onChange({ ...value, version: e.target.value })}
              placeholder="2024-01-01"
              className="h-8 w-48 text-sm"
            />
          </div>

          {value.statement.map((stmt, i) => (
            <StatementCard
              key={i}
              index={i}
              statement={stmt}
              onChange={(s) => updateStatement(i, s)}
              onRemove={() => removeStatement(i)}
              canRemove={value.statement.length > 1}
            />
          ))}

          <Button type="button" variant="outline" onClick={addStatement}>
            <Plus className="mr-1 h-4 w-4" />
            Add Statement
          </Button>
        </TabsContent>

        <TabsContent value="json">
          <JsonEditor
            ref={editorRef}
            defaultValue={JSON.stringify(value, null, 2)}
            resetKey={jsonResetKey}
            onChange={(raw) => {
              try {
                const parsed = JSON.parse(raw);
                onChange(parsed);
              } catch {
                // Invalid JSON — don't update form state
              }
            }}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
