#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

kubectl apply --server-side -k platform/argocd
kubectl apply --server-side -k platform/rollouts
kubectl apply --server-side -k platform/prometheus

kubectl -n argocd rollout status deploy/argocd-server --timeout=5m
kubectl -n argo-rollouts rollout status deploy/argo-rollouts --timeout=5m

kubectl apply -f platform/argocd/apps/

echo
echo "argocd admin password:"
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d
echo
