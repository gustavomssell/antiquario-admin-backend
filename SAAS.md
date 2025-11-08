DOCUMENTAÇÃO - ESTRUTURA DO BANCO DE DADOS
SaaS para Antiquário
🎯 VISÃO GERAL DO SISTEMA
Objetivo:
Sistema completo para gestão de antiquário com:

Backend API: Gestão administrativa completa
Dashboard Admin: Interface para administradores e vendedores
Site Público: Catálogo para compradores + sistema de leilão
Tipos de Usuários:
Administradores: Controle total do sistema
Vendedores: Gestão de produtos, clientes, vendas
Compradores: Visualização, compras, leilões
🗄️ ESTRUTURA DETALHADA DO BANCO
👥 1. GESTÃO DE USUÁRIOS E PERMISSÕES
users (já existe no template)
sql
Copiar

- id (PK)
- name
- email (unique)
- password
- role_id (FK)
- created_at, updated_at, deleted_at
roles
sql
Copiar

- id (PK)
- name (admin, vendedor, comprador)
- description
- permissions (JSON ou tabela separada)
- created_at, updated_at
user_profiles
sql
Copiar

- id (PK)
- user_id (FK)
- phone, address (JSON)
- document_type, document_number
- birth_date
- created_at, updated_at
🏷️ 2. SISTEMA DE CATEGORIZAÇÃO
categories (Hierárquica - 3 níveis)
sql
Copiar

- id (PK)
- name
- description
- parent_id (FK self-reference)
- level (1, 2, 3)
- image_url
- active (boolean)
- created_by (FK users)
- created_at, updated_at
Exemplo de Hierarquia:

Nível 1: Móveis, Arte, Joias, Decoração
Nível 2: Assentos, Mesas, Armários (para Móveis)
Nível 3: Cadeira de Jantar, Poltrona, Banco (para Assentos)
periods (Períodos Históricos)
sql
Copiar

- id (PK)
- name (Barroco, Art Déco, etc.)
- description
- start_year, end_year
- created_at, updated_at
styles (Estilos por Origem)
sql
Copiar

- id (PK)
- name (Colonial Brasileiro, Francês, etc.)
- description
- period_id (FK)
- origin_country
- created_at, updated_at
materials
sql
Copiar

- id (PK)
- name (Madeira de Lei, Bronze, etc.)
- description
- category (madeira, metal, tecido, etc.)
- created_at, updated_at
tags (Tags Livres)
sql
Copiar

- id (PK)
- name (raro, restaurado, assinado, etc.)
- color (para UI)
- description
- active (boolean)
- created_at, updated_at
product_tags (Many-to-Many)
sql
Copiar

- id (PK)
- product_id (FK)
- tag_id (FK)
📦 3. PRODUTOS E CONJUNTOS
product_sets (Conjuntos de Produtos)
sql
Copiar

- id (PK)
- set_code (ANT-2024-SET-001)
- name, description
- total_pieces
- can_sell_separately (boolean)
- created_by (FK users)
- created_at, updated_at
products (Tabela Principal)
sql
Copiar

- id (PK)
- code (ANT-2024-0001 ou ANT-2024-SET-001-A) (unique)
- qr_code_url
- product_set_id (FK product_sets) (nullable)
- set_position (A, B, C...) (nullable)
- name, description
- category_id (FK), period_id (FK), style_id (FK)
- dimensions (JSON: altura, largura, profundidade)
- weight, condition_rating (1-10)
- acquisition_type (purchase/consignment)
- acquisition_date, acquisition_price
- consignment_percentage, consignment_deadline
- supplier_id (FK)
- estimated_value, selling_price, commission_rate
- status (ENUM: available, sold, auction, reserved, restoration, 
         evaluation, consignment, damaged, exhibition, returned)
- provenance_story (TEXT), historical_notes (TEXT)
- is_set_item (boolean), available_quantity
- created_by (FK), updated_by (FK)
- created_at, updated_at, deleted_at
Códigos de Exemplo:

Produto único: ANT-2024-0001
Conjunto: ANT-2024-SET-001
Itens do conjunto: ANT-2024-SET-001-A, ANT-2024-SET-001-B
product_materials (Many-to-Many)
sql
Copiar

- id (PK)
- product_id (FK)
- material_id (FK)
- percentage (para peças com múltiplos materiais)
- notes
📱 4. MÍDIA E DOCUMENTOS
product_media (Fotos e Vídeos)
sql
Copiar

- id (PK)
- product_id (FK)
- media_type (image/video)
- media_url
- alt_text, is_primary (boolean)
- order_position, file_size, duration (para vídeos)
- created_at, updated_at
Limite: 10 mídias por produto (fotos + vídeos)

product_documents (PDFs)
sql
Copiar

- id (PK)
- product_id (FK)
- document_type (certificate/appraisal/invoice/provenance)
- document_url, title, description, file_size
- created_at, updated_at
Limite: 3 PDFs por produto

🏪 5. GESTÃO COMERCIAL
suppliers (Fornecedores)
sql
Copiar

- id (PK)
- name, type (pessoa_fisica/pessoa_juridica)
- contact_info (JSON: telefones, emails)
- address (JSON), document_number
- notes
- created_at, updated_at
customers (Compradores)
sql
Copiar

- id (PK)
- user_id (FK)
- customer_type (individual/empresa)
- preferences (JSON)
- purchase_history_summary (JSON)
- credit_limit, notes
- created_at, updated_at
acquisitions (Entrada de Peças)
sql
Copiar

- id (PK)
- supplier_id (FK)
- acquisition_date, total_value, payment_method
- notes, created_by (FK)
- created_at, updated_at
acquisition_items
sql
Copiar

- id (PK)
- acquisition_id (FK)
- product_id (FK)
- unit_price, quantity, total_price
💰 6. VENDAS E TRANSAÇÕES
sales
sql
Copiar

- id (PK)
- customer_id (FK)
- sale_date, sale_type (direct/auction/set)
- total_amount, payment_method, payment_status
- delivery_address (JSON), delivery_date
- notes, created_by (FK)
- created_at, updated_at
sale_items
sql
Copiar

- id (PK)
- sale_id (FK)
- product_id (FK)
- unit_price, quantity, total_price
- commission_percentage, commission_amount, commission_paid
payments
sql
Copiar

- id (PK)
- sale_id (FK)
- amount, payment_date, payment_method
- status, reference
- created_at, updated_at
💼 7. SISTEMA DE COMISSÕES
commission_rules
sql
Copiar

- id (PK)
- user_id (FK) (vendedor)
- product_category_id (FK) (nullable - se específico por categoria)
- commission_percentage
- is_default (boolean), active (boolean)
- created_by (FK), created_at, updated_at
🏦 8. CONSIGNAÇÃO
consignment_settings (Configurações Globais)
sql
Copiar

- id (PK)
- default_percentage
- default_deadline_days
- created_by (FK)
- created_at, updated_at
consignment_returns (Devoluções)
sql
Copiar

- id (PK)
- product_id (FK)
- return_date, return_reason
- condition_on_return (TEXT)
- returned_by (FK users), notes
- created_at, updated_at
🔒 9. SISTEMA DE RESERVAS
reservation_settings (Configurações)
sql
Copiar

- id (PK)
- default_reservation_days
- requires_deposit (boolean), deposit_percentage
- auto_qualification_enabled (boolean)
- min_purchases_for_auto, min_total_spent
- good_reputation_threshold
- created_by (FK), updated_at
customer_qualifications (Habilitação para Reservas)
sql
Copiar

- id (PK)
- customer_id (FK)
- qualification_type (manual/auto/vip/premium)
- is_active (boolean)
- qualification_date, qualified_by (FK users)
- auto_qualification_reason (purchases_count/total_spent/reputation)
- notes, created_by (FK)
- created_at, updated_at
customer_reputation (Reputação dos Clientes)
sql
Copiar

- id (PK)
- customer_id (FK)
- total_purchases, total_spent, avg_rating
- late_payments_count, cancelled_reservations_count
- reputation_score (calculado automaticamente)
- last_calculated
- created_at, updated_at
reservations
sql
Copiar

- id (PK)
- product_id (FK), customer_id (FK)
- reservation_date, expiry_date
- custom_reservation_days (override do padrão)
- requires_deposit (boolean), deposit_percentage
- deposit_amount, deposit_paid (boolean), deposit_date
- status (active/expired/converted/cancelled)
- cancellation_reason, notes
- created_by (FK)
- created_at, updated_at
🔍 10. AVALIAÇÃO E AUTENTICAÇÃO
appraisals (Avaliações)
sql
Copiar

- id (PK)
- product_id (FK)
- appraiser_name, appraisal_date
- estimated_value, condition_assessment
- authenticity_rating, notes, document_url
- created_at, updated_at
certificates (Certificados)
sql
Copiar

- id (PK)
- product_id (FK)
- certificate_type, issuer
- issue_date, expiry_date
- certificate_number, document_url
- verified (boolean)
- created_at, updated_at
🎯 11. SISTEMA DE LEILÃO (Futuro)
auctions
sql
Copiar

- id (PK)
- title, description
- start_date, end_date
- status, created_by (FK)
- created_at, updated_at
auction_items
sql
Copiar

- id (PK)
- auction_id (FK), product_id (FK)
- starting_bid, reserve_price, current_bid
bids
sql
Copiar

- id (PK)
- auction_item_id (FK), bidder_id (FK)
- bid_amount, bid_date, status
🔄 FLUXOS PRINCIPAIS
📥 Entrada de Peças:
Recebimento (fornecedor/consignação)
Catalogação inicial (status: evaluation)
Avaliação por especialista
Certificação/autenticação
Precificação
Disponibilização (status: available)
📤 Saída de Peças:
Venda direta ou leilão
Emissão de certificados
Documentação de proveniência
Entrega/logística
🔒 Qualificação para Reservas:
Manual: Administrador habilita/desabilita
Automática: Baseada em critérios configuráveis:
Número mínimo de compras
Valor total gasto
Score de reputação
⚙️ CONFIGURAÇÕES ADMINISTRATIVAS
🎛️ Configurações Globais:
Percentual padrão de consignação
Prazo padrão de reservas (dias)
Critérios de auto-qualificação para reservas
Exigência de depósito para reservas
Percentual de comissão por categoria/vendedor
📊 Status dos Produtos:
available - Disponível para venda
sold - Vendido
auction - Em leilão
reserved - Reservado por cliente
restoration - Em restauração
evaluation - Em avaliação/autenticação
consignment - Consignado (aguardando decisão)
damaged - Danificado (fora de linha)
exhibition - Em exposição (não disponível)
returned - Devolvido (consignação)
🏷️ Tags Sugeridas:
Condição: "restaurado", "original", "assinado"
Raridade: "raro", "único", "edição limitada"
Destaque: "peça do mês", "promoção", "novidade"
Origem: "importado", "nacional", "regional"
🔐 CONSTRAINTS E VALIDAÇÕES
📏 Limites:
Mídia: Máximo 10 arquivos por produto (fotos + vídeos)
Documentos: Máximo 3 PDFs por produto
Códigos: Únicos no sistema
QR Codes: Gerados automaticamente
🔗 Relacionamentos Principais:
Produtos podem pertencer a conjuntos (1:N)
Produtos têm múltiplas mídias (1:N)
Produtos têm múltiplos materiais (N:N)
Produtos têm múltiplas tags (N:N)
Clientes podem ter múltiplas reservas (1:N)
Vendas podem ter múltiplos itens (1:N)
🚀 FUNCIONALIDADES FUTURAS
📱 Arquivo 3D:
Suporte para arquivos GLB/GLTF
Visualização 3D no catálogo
Realidade aumentada
🤖 Automações:
Expiração automática de reservas
Cálculo automático de reputação
Notificações de prazos de consignação
Geração automática de relatórios
Documentação criada para o SaaS de Antiquário - Versão 1.0



# 📚 DOCUMENTAÇÃO - ESTRUTURA DO BANCO DE DADOS
**SaaS para Antiquário — Versão 1.0**

## 🎯 Visão Geral do Sistema
**Objetivo:** Sistema completo para gestão de antiquário com:

- **Backend API**: Gestão administrativa e integrações
- **Dashboard Administrativo**: Para administradores e vendedores
- **Site Público**: Catálogo e sistema de leilão para compradores

**Tipos de Usuários:**

| Tipo | Permissões principais |
|-----|------------------------|
| Administradores | Controle total |
| Vendedores | Gestão de produtos e vendas |
| Compradores | Visualização, compras e participação em leilões |

## 🗄️ Estrutura Detalhada do Banco

### 1. Gestão de Usuários e Permissões

**Tabela: `users`**
```
id (PK)
name
email (unique)
password
role_id (FK)
created_at
updated_at
deleted_at
```

**Tabela: `roles`**
```
id (PK)
name (admin, vendedor, comprador)
description
permissions (JSON ou tabela relacionada)
created_at
updated_at
```

**Tabela: `user_profiles`**
```
id (PK)
user_id (FK)
phone
address (JSON)
document_type
document_number
birth_date
created_at
updated_at
```

### 2. Sistema de Categorização

**Tabela: `categories`**
```
id (PK)
name
description
parent_id (FK self-reference)
level (1, 2, 3)
image_url
active (boolean)
created_by (FK users)
created_at
updated_at
```

**Tabela: `periods`**
```
id (PK)
name
description
start_year
end_year
created_at
updated_at
```

**Tabela: `styles`**
```
id (PK)
name
description
period_id (FK)
origin_country
created_at
updated_at
```

**Tabela: `materials`**
```
id (PK)
name
description
category
created_at
updated_at
```

**Tabela: `tags`**
```
id (PK)
name
color
description
active (boolean)
created_at
updated_at
```

**Tabela Relacional: `product_tags`**
```
id (PK)
product_id (FK)
tag_id (FK)
```

### 3. Produtos e Conjuntos

**Tabela: `product_sets`**
```
id (PK)
set_code
name
description
total_pieces
can_sell_separately (boolean)
created_by (FK users)
created_at
updated_at
```

**Tabela: `products`**
```
id (PK)
code (unique)
qr_code_url
product_set_id (FK, nullable)
set_position (A, B, C…)
name
description
category_id (FK)
period_id (FK)
style_id (FK)
dimensions (JSON)
weight
condition_rating (1-10)
acquisition_type (purchase/consignment)
acquisition_date
acquisition_price
consignment_percentage
consignment_deadline
supplier_id (FK)
estimated_value
selling_price
commission_rate
status
provenance_story (TEXT)
historical_notes (TEXT)
is_set_item (boolean)
available_quantity
created_by (FK)
updated_by (FK)
created_at
updated_at
deleted_at
```

**Tabela: `product_materials`**
```
id (PK)
product_id (FK)
material_id (FK)
percentage
notes
```

### 4. Mídia e Documentos

**Tabela: `product_media`**
```
id (PK)
product_id (FK)
media_type (image/video)
media_url
alt_text
is_primary (boolean)
order_position
file_size
duration
created_at
updated_at
```

**Tabela: `product_documents`**
```
id (PK)
product_id (FK)
document_type
document_url
title
description
file_size
created_at
updated_at
```

### 5. Gestão Comercial

**Tabela: `suppliers`**
```
id (PK)
name
type (pessoa_fisica/pessoa_juridica)
contact_info (JSON)
address (JSON)
document_number
notes
created_at
updated_at
```

**Tabela: `customers`**
```
id (PK)
user_id (FK)
customer_type
preferences (JSON)
purchase_history_summary (JSON)
credit_limit
notes
created_at
updated_at
```

**Tabela: `acquisitions`**
```
id (PK)
supplier_id (FK)
acquisition_date
total_value
payment_method
notes
created_by (FK)
created_at
updated_at
```

**Tabela: `acquisition_items`**
```
id (PK)
acquisition_id (FK)
product_id (FK)
unit_price
quantity
total_price
```

### 6. Vendas e Transações

**Tabela: `sales`**
```
id (PK)
customer_id (FK)
sale_date
sale_type
total_amount
payment_method
payment_status
delivery_address (JSON)
delivery_date
notes
created_by (FK)
created_at
updated_at
```

**Tabela: `sale_items`**
```
id (PK)
sale_id (FK)
product_id (FK)
unit_price
quantity
total_price
commission_percentage
commission_amount
commission_paid
```

**Tabela: `payments`**
```
id (PK)
sale_id (FK)
amount
payment_date
payment_method
status
reference
created_at
updated_at
```

### 7. Sistema de Comissões

**Tabela: `commission_rules`**
```
id (PK)
user_id (FK)
product_category_id (FK, nullable)
commission_percentage
is_default (boolean)
active (boolean)
created_by (FK)
created_at
updated_at
```

### 8. Consignação

**Tabela: `consignment_settings`**
```
id (PK)
default_percentage
default_deadline_days
created_by (FK)
created_at
updated_at
```

**Tabela: `consignment_returns`**
```
id (PK)
product_id (FK)
return_date
return_reason
condition_on_return
returned_by (FK)
notes
created_at
updated_at
```

### 9. Sistema de Reservas

**Tabela: `reservation_settings`**
```
id (PK)
default_reservation_days
requires_deposit (boolean)
deposit_percentage
auto_qualification_enabled (boolean)
min_purchases_for_auto
min_total_spent
good_reputation_threshold
created_by (FK)
updated_at
```

**Tabela: `customer_qualifications`**
```
id (PK)
customer_id (FK)
qualification_type
is_active (boolean)
qualification_date
qualified_by (FK)
auto_qualification_reason
notes
created_by (FK)
created_at
updated_at
```

**Tabela: `customer_reputation`**
```
id (PK)
customer_id (FK)
total_purchases
total_spent
avg_rating
late_payments_count
cancelled_reservations_count
reputation_score
last_calculated
created_at
updated_at
```

**Tabela: `reservations`**
```
id (PK)
product_id (FK)
customer_id (FK)
reservation_date
expiry_date
requires_deposit
deposit_amount
status
cancellation_reason
notes
created_by (FK)
created_at
updated_at
```

### 10. Avaliação e Autenticação

**Tabela: `appraisals`**
```
id (PK)
product_id (FK)
appraiser_name
appraisal_date
estimated_value
condition_assessment
authenticity_rating
notes
document_url
created_at
updated_at
```

**Tabela: `certificates`**
```
id (PK)
product_id (FK)
certificate_type
issuer
issue_date
expiry_date
certificate_number
document_url
verified (boolean)
created_at
updated_at
```

### 11. Sistema de Leilão (Futuro)

**Tabela: `auctions`**
```
id (PK)
title
description
start_date
end_date
status
created_by (FK)
created_at
updated_at
```

**Tabela: `auction_items`**
```
id (PK)
auction_id (FK)
product_id (FK)
starting_bid
reserve_price
current_bid
```

**Tabela: `bids`**
```
id (PK)
auction_item_id (FK)
bidder_id (FK)
bid_amount
bid_date
status
```