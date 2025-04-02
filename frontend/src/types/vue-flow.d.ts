declare module '@vue-flow/core' {
  export interface Node {
    id: string;
    type?: string;
    position: { x: number; y: number };
    data?: any;
    style?: any;
    [key: string]: any;
  }

  export interface Edge {
    id: string;
    source: string;
    target: string;
    label?: string;
    type?: string;
    style?: any;
    labelStyle?: any;
    labelBgStyle?: any;
    animated?: boolean;
    markerEnd?: any;
    [key: string]: any;
  }

  export enum MarkerType {
    Arrow = 'arrow',
    ArrowClosed = 'arrowclosed'
  }

  export interface NodeMouseEvent {
    node: Node;
    event: MouseEvent;
  }

  export const VueFlow: any;
  export function useVueFlow(options: any): any;
  export function useZoomPanHelper(): any;
  export const Panel: any;
}

declare module '@vue-flow/background';
declare module '@vue-flow/controls';
declare module '@vue-flow/minimap'; 