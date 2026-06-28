import React, { useRef, useState, useEffect } from "react";
import { UploadCloud, X, File as FileIcon } from "lucide-react";
import { useErrorHandler } from "../../hooks/useErrorHandler";

import {
  DndContext,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";

import {
  SortableContext,
  useSortable,
  arrayMove,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";

import { CSS } from "@dnd-kit/utilities";
import type { FileInfo } from "../../types/common";
import CommonService from "../../services/common/common.service";
import { apiUrl } from "../../services/api";
import { Box, Flex, Text, Card, IconButton } from "@radix-ui/themes";

function SortableItem({ id, children }: any) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    cursor: "grab",
    position: "relative",
    zIndex: isDragging ? 9999 : "auto",
  };

  return (
    <Box ref={setNodeRef} style={style} {...attributes} {...listeners}>
      {children}
    </Box>
  );
}

export interface CUploadProps {
  multiple?: boolean;
  accept?: string;
  maxSizeMB?: number;
  maxFiles?: number;
  existingFiles?: FileInfo[];
  onFilesChange?: (data: {
    newFiles: FileInfo[];
    keepFiles: FileInfo[];
    order: FileInfo[];
  }) => void;
}

const CUpload: React.FC<CUploadProps> = ({
  multiple = false,
  accept = "image/*",
  maxSizeMB = 5,
  maxFiles = 10,
  existingFiles = [],
  onFilesChange,
}) => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const { showError } = useErrorHandler();

  const [uploaded, setUploaded] = useState<FileInfo[]>([]);
  const [progress, setProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);

  useEffect(() => {
    setUploaded(existingFiles);
  }, [existingFiles]);

  const emitChange = (newFiles: FileInfo[], keepFiles: FileInfo[]) => {
    onFilesChange?.({
      newFiles,
      keepFiles,
      order: [...keepFiles, ...newFiles],
    });
  };

  const handleSelectFiles = async (list: FileList | null) => {
    if (!list) return;

    const files = Array.from(list);

    if (uploaded.length + files.length > maxFiles) {
      showError(`Max ${maxFiles} files allowed`);
      return;
    }

    const valid = files.filter((f) => f.size / 1024 / 1024 <= maxSizeMB);
    if (valid.length !== files.length)
      showError(`Some files exceeded ${maxSizeMB} MB`);
    if (valid.length === 0) return;

    try {
      setIsUploading(true);
      const formData = new FormData();
      valid.forEach((f) => formData.append("files", f));

      const newUploaded = await CommonService.uploadWithProgress(
        formData,
        (p) => setProgress(p)
      );

      setUploaded((prev) => {
        const updated = [...prev, ...newUploaded];
        const keepFiles = updated.filter((f) =>
          existingFiles.some((ex) => ex.ID === f.ID)
        );
        const newFiles = updated.filter(
          (f) => !existingFiles.some((ex) => ex.ID === f.ID)
        );
        emitChange(newFiles, keepFiles);
        return updated;
      });
    } catch {
      showError("Upload failed");
    } finally {
      setIsUploading(false);
      setProgress(0);
    }
  };

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } })
  );

  const handleDragEnd = (result: any) => {
    if (!result.over) return;
    const oldIndex = uploaded.findIndex((i) => i.ID === result.active.id);
    const newIndex = uploaded.findIndex((i) => i.ID === result.over.id);
    if (oldIndex === newIndex) return;

    setUploaded((prev) => {
      const updated = arrayMove(prev, oldIndex, newIndex);
      const keepFiles = updated.filter((f) =>
        existingFiles.some((ex) => ex.ID === f.ID)
      );
      const newFiles = updated.filter(
        (f) => !existingFiles.some((ex) => ex.ID === f.ID)
      );
      emitChange(newFiles, keepFiles);
      return updated;
    });
  };

  const handleRemove = (file: FileInfo) => {
    const updated = uploaded.filter((f) => f.ID !== file.ID);
    const keepFiles = updated.filter((f) =>
      existingFiles.some((ex) => ex.ID === f.ID)
    );
    const newFiles = updated.filter(
      (f) => !existingFiles.some((ex) => ex.ID === f.ID)
    );
    setUploaded(updated);
    emitChange(newFiles, keepFiles);
  };

  const isImage = (file: FileInfo) => file.Mime?.startsWith("image/");

  return (
    <Box style={{ overflow: "visible" }}>
      <Box
        style={{
          border: "1px dashed var(--gray-6)",
          cursor: "pointer",
          padding: "20px",
          borderRadius: "2px",
          marginBottom: "12px",
          marginTop: "6px",
        }}
        onClick={() => inputRef.current?.click()}
        onDrop={(e) => {
          e.preventDefault();
          handleSelectFiles(e.dataTransfer.files);
        }}
        onDragOver={(e) => e.preventDefault()}
      >
        <Flex align="center" gap="3" direction="column">
          <UploadCloud size={26} />
          <Text style={{ fontSize: "12px" }}>
            Drag & drop files or <strong>browse</strong>
          </Text>
          <input
            ref={inputRef}
            hidden
            type="file"
            accept={accept}
            multiple={multiple}
            onChange={(e) => handleSelectFiles(e.target.files)}
          />
        </Flex>
      </Box>

      {isUploading && (
        <Box mt="2" height="6px" width="100%">
          <Box height="100%" style={{ width: `${progress}%` }} />
        </Box>
      )}

      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
      >
        <SortableContext
          items={uploaded.map((f) => f.ID)}
          strategy={verticalListSortingStrategy}
        >
          <Flex direction="column" gap="2" mt="3">
            {uploaded.map((file) => (
              <SortableItem key={file.ID} id={file.ID}>
                <Card
                  style={{
                    padding: "8px",
                    ["--card-border-radius" as any]: "3px",
                  }}
                >
                  <Flex align="center" gap="3" p="1">
                    <Box
                      width="60px"
                      height="60px"
                      style={{ overflow: "hidden", borderRadius: 3 }}
                    >
                      {isImage(file) ? (
                        <img
                          src={apiUrl(file.URL)}
                          style={{
                            width: "100%",
                            height: "100%",
                            objectFit: "cover",
                          }}
                        />
                      ) : (
                        <FileIcon size={20} />
                      )}
                    </Box>

                    <Flex direction="column" flexGrow="1">
                      <Text weight="medium" style={{ fontSize: "13px" }}>
                        {file.Name}
                      </Text>
                      <Text size="1">
                        {(file.Size / 1024 / 1024).toFixed(2)} MB
                      </Text>
                    </Flex>

                    <IconButton
                      style={{
                        padding: "5px",
                        marginRight: "10px",
                        cursor: "pointer",
                      }}
                      variant="ghost"
                      radius="full"
                      onClick={() => handleRemove(file)}
                    >
                      <X size={16} />
                    </IconButton>
                  </Flex>
                </Card>
              </SortableItem>
            ))}
          </Flex>
        </SortableContext>
      </DndContext>
    </Box>
  );
};

export default CUpload;
