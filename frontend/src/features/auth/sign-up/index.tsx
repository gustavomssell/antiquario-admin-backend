import { Link } from '@tanstack/react-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { AuthLayout } from '../auth-layout'
import { SignUpForm } from './components/sign-up-form'

export function SignUp() {
  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader>
          <CardTitle className='text-lg tracking-tight'>
            Criar conta de Administrador
          </CardTitle>
          <CardDescription>
            Preencha seus dados para receber acesso interno. <br />
            Já possui uma conta?{' '}
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              Fazer Login
            </Link>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SignUpForm />
        </CardContent>
        <CardFooter>
          <p className='text-muted-foreground px-8 text-center text-sm'>
            Ao criar uma conta, você concorda com nossos{' '}
            <a
              href='/terms'
              className='hover:text-primary underline underline-offset-4'
            >
              Termos de Serviço
            </a>{' '}
            e{' '}
            <a
              href='/privacy'
              className='hover:text-primary underline underline-offset-4'
            >
              Política de Privacidade
            </a>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
