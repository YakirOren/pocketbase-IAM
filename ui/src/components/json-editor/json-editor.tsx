import { useRef, useEffect, forwardRef, useImperativeHandle } from "react";
import { EditorState } from "@codemirror/state";
import type { Extension } from "@codemirror/state";
import {
  EditorView,
  keymap,
  lineNumbers,
  highlightActiveLine,
  highlightSpecialChars,
  drawSelection,
} from "@codemirror/view";
import {
  defaultKeymap,
  history,
  historyKeymap,
  indentWithTab,
} from "@codemirror/commands";
import { json } from "@codemirror/lang-json";
import { linter, lintGutter } from "@codemirror/lint";
import type { Diagnostic } from "@codemirror/lint";
import {
  syntaxHighlighting,
  defaultHighlightStyle,
  indentOnInput,
  bracketMatching,
  foldGutter,
  foldKeymap,
} from "@codemirror/language";
import { closeBrackets, closeBracketsKeymap } from "@codemirror/autocomplete";
import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
import { cn } from "@/lib/utils";

const jsonLinter = linter((view) => {
  const diagnostics: Diagnostic[] = [];
  const doc = view.state.doc.toString();
  if (doc.trim().length === 0) return diagnostics;

  try {
    JSON.parse(doc);
  } catch (e) {
    if (e instanceof SyntaxError) {
      const posMatch = e.message.match(/position\s+(\d+)/i);
      const pos = posMatch ? Number(posMatch[1]) : 0;
      diagnostics.push({
        from: Math.min(pos, doc.length),
        to: Math.min(pos, doc.length),
        severity: "error",
        message: e.message,
      });
    }
  }
  return diagnostics;
});

export interface JsonEditorHandle {
  getValue: () => string;
  setValue: (value: string) => void;
}

interface JsonEditorProps {
  defaultValue?: string;
  onChange?: (value: string) => void;
  readOnly?: boolean;
  className?: string;
  resetKey?: string | number;
}

export const JsonEditor = forwardRef<JsonEditorHandle, JsonEditorProps>(
  function JsonEditor(
    { defaultValue = "", onChange, readOnly = false, className, resetKey },
    ref
  ) {
    const containerRef = useRef<HTMLDivElement>(null);
    const viewRef = useRef<EditorView | null>(null);
    const onChangeRef = useRef(onChange);
    onChangeRef.current = onChange;

    useImperativeHandle(ref, () => ({
      getValue() {
        return viewRef.current?.state.doc.toString() ?? "";
      },
      setValue(value: string) {
        const view = viewRef.current;
        if (!view) return;
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: value },
        });
      },
    }));

    useEffect(() => {
      if (!containerRef.current) return;
      viewRef.current?.destroy();

      const extensions: Extension[] = [
        lineNumbers(),
        highlightActiveLine(),
        highlightSpecialChars(),
        drawSelection(),
        indentOnInput(),
        bracketMatching(),
        closeBrackets(),
        foldGutter(),
        history(),
        highlightSelectionMatches(),
        syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
        keymap.of([
          ...closeBracketsKeymap,
          ...defaultKeymap,
          ...historyKeymap,
          ...foldKeymap,
          ...searchKeymap,
          indentWithTab,
        ]),
        json(),
        jsonLinter,
        lintGutter(),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            onChangeRef.current?.(update.state.doc.toString());
          }
        }),
        EditorView.theme({
          "&": { minHeight: "200px", maxHeight: "500px" },
          ".cm-scroller": { overflow: "auto" },
        }),
      ];

      if (readOnly) {
        extensions.push(EditorState.readOnly.of(true));
        extensions.push(EditorView.editable.of(false));
      }

      const view = new EditorView({
        state: EditorState.create({ doc: defaultValue, extensions }),
        parent: containerRef.current,
      });
      viewRef.current = view;

      return () => {
        view.destroy();
        viewRef.current = null;
      };
    }, [resetKey, readOnly]);

    return (
      <div
        ref={containerRef}
        className={cn(
          "overflow-hidden rounded-md border",
          className
        )}
      />
    );
  }
);
