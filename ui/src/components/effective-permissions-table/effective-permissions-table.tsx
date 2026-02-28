import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import type { EffectivePermission } from "@/lib/permissions";

interface EffectivePermissionsTableProps {
  permissions: EffectivePermission[];
  showSource?: boolean;
}

export function EffectivePermissionsTable({
  permissions,
  showSource = false,
}: EffectivePermissionsTableProps) {
  if (permissions.length === 0) {
    return (
      <p className="py-4 text-sm text-muted-foreground">
        No permissions to display.
      </p>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Action</TableHead>
          <TableHead>Effect</TableHead>
          {showSource && <TableHead>Source</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {permissions.map((p) => (
          <TableRow key={p.action}>
            <TableCell className="font-mono text-sm">{p.action}</TableCell>
            <TableCell>
              <Badge
                variant={p.effect === "Allow" ? "default" : "destructive"}
              >
                {p.effect}
              </Badge>
            </TableCell>
            {showSource && (
              <TableCell className="text-sm text-muted-foreground">
                {p.source}
              </TableCell>
            )}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
