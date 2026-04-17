import { HttpErrorResponse } from '@angular/common/http';

type ErrorPayload = {
  message?: string;
  error?: string;
  status?: string;
};

function readPayload(error: unknown): ErrorPayload {
  if (!(error instanceof HttpErrorResponse)) {
    return {};
  }

  if (typeof error.error === 'object' && error.error !== null) {
    return error.error as ErrorPayload;
  }

  return {};
}

function includesAny(value: string, patterns: string[]): boolean {
  const normalized = value.toLowerCase();
  return patterns.some((pattern) => normalized.includes(pattern.toLowerCase()));
}

export function mapProductError(error: unknown, fallback: string): string {
  if (!(error instanceof HttpErrorResponse)) {
    return fallback;
  }

  const payload = readPayload(error);
  const combined = `${payload.message ?? ''} ${payload.error ?? ''}`;

  if (error.status === 0) {
    return 'Problema no servico de produtos, tente novamente mais tarde.';
  }

  if (error.status === 400) {
    if (includesAny(combined, ['codigo and descricao are required'])) {
      return 'Informe codigo e descricao do produto.';
    }

    if (includesAny(combined, ['quantidade must be greater than zero'])) {
      return 'Informe uma quantidade maior que zero.';
    }

    return payload.message ?? fallback;
  }

  if (error.status === 404 && includesAny(combined, ['product not found'])) {
    return 'Produto nao encontrado.';
  }

  if (error.status === 409 && includesAny(combined, ['linked to existing invoices', 'product in use'])) {
    return 'Nao e possivel remover o produto porque ele ja esta vinculado a notas fiscais.';
  }

  if (error.status === 409 || includesAny(combined, ['duplicate key', 'already exists', 'unique'])) {
    return 'Ja existe um produto cadastrado com esse codigo.';
  }

  if (error.status >= 500) {
    return 'Problema no servico de produtos, tente novamente mais tarde.';
  }

  return payload.message ?? fallback;
}

export function mapInvoiceCreationError(error: unknown, fallback: string): string {
  if (!(error instanceof HttpErrorResponse)) {
    return fallback;
  }

  const payload = readPayload(error);
  const combined = `${payload.message ?? ''} ${payload.error ?? ''}`;

  if (error.status === 0) {
    return 'Problema no servico de faturamento, tente novamente mais tarde.';
  }

  if (error.status === 400) {
    if (includesAny(combined, ['itens is required'])) {
      return 'Adicione pelo menos um item na nota fiscal.';
    }

    if (includesAny(combined, ['codigo_produto is required'])) {
      return 'Todos os itens da nota fiscal precisam ter codigo do produto.';
    }

    if (includesAny(combined, ['quantidade must be greater than zero'])) {
      return 'Todas as quantidades da nota fiscal devem ser maiores que zero.';
    }

    if (includesAny(combined, ['codigo_produto must be unique within itens'])) {
      return 'Nao repita o mesmo produto na nota fiscal. Ajuste a quantidade no proprio item.';
    }

    return payload.message ?? fallback;
  }

  if (error.status === 404 && includesAny(combined, ['product not found'])) {
    return 'Um dos produtos informados nao foi encontrado no estoque.';
  }

  if (error.status >= 500) {
    if (includesAny(combined, ['estoque', 'product'])) {
      return 'Problema no servico de produtos, tente novamente mais tarde.';
    }
    return 'Problema no servico de faturamento, tente novamente mais tarde.';
  }

  return payload.message ?? fallback;
}

export function mapInvoicePrintError(error: unknown, fallback: string): string {
  if (!(error instanceof HttpErrorResponse)) {
    return fallback;
  }

  const payload = readPayload(error);
  const combined = `${payload.message ?? ''} ${payload.error ?? ''}`;

  if (error.status === 0) {
    return 'Problema no servico de faturamento, tente novamente mais tarde.';
  }

  if (error.status === 404) {
    if (includesAny(combined, ['nota fiscal not found'])) {
      return 'A nota fiscal informada nao foi encontrada.';
    }

    if (includesAny(combined, ['product not found'])) {
      return 'Um produto vinculado a essa nota fiscal nao foi encontrado no estoque.';
    }
  }

  if (error.status === 409) {
    if (includesAny(combined, ['already closed'])) {
      return 'Essa nota fiscal ja foi impressa.';
    }

    if (includesAny(combined, ['insufficient stock'])) {
      return 'Um ou mais produtos nao possuem saldo suficiente para imprimir a nota fiscal.';
    }
  }

  if (error.status >= 500) {
    if (includesAny(combined, ['estoque', 'product'])) {
      return 'Problema no servico de produtos, tente novamente mais tarde.';
    }
    return 'Problema no servico de faturamento, tente novamente mais tarde.';
  }

  return payload.message ?? fallback;
}

export function mapInvoiceDetailError(error: unknown, fallback: string): string {
  if (!(error instanceof HttpErrorResponse)) {
    return fallback;
  }

  const payload = readPayload(error);
  const combined = `${payload.message ?? ''} ${payload.error ?? ''}`;

  if (error.status === 0 || error.status >= 500) {
    return 'Problema no servico de faturamento, tente novamente mais tarde.';
  }

  if (error.status === 404 && includesAny(combined, ['nota fiscal not found'])) {
    return 'A nota fiscal informada nao foi encontrada.';
  }

  return payload.message ?? fallback;
}
